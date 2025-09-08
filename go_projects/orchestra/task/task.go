package task

import (
	"context"
	"io"
	"log"
	"math"
	"os"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
)

//here we define the states a task can be assigned
type State int

//task options
const (
	Pending State = iota
	Scheduled
	Runnig
	Completed
	Failed
)

//Task attributes and Methods structure type
type Task struct {
	ID uuid.UUID
	Name string
	State State
	Image string //represents the docker image,will be used by the scheduler to find a cluster capable of running a task
	ContainerId string
	Memory int //resource metric
	Disk int //resource metric
	ExposedPorts nat.PortSet //allocation of proper ports for the task
	PortBindings map[string]string //allcation of proper ports foor the task
	Env string
	RestartPolicy string //accepts values such as "always", "unless-stopped", "on-failure"
	StartTime  time.Time
	Finish time.Time
}


//event that can be triggered by the task
type TaskEvent struct {
	ID uuid.UUID //unique identifier of the a task event
	State State //rep the progression of the state of th task
	Timestamp time.Time //record the time the event was requested
	Task Task //the task that requested the taskevent
}

//Configuration for the tasks 
type Config struct {
	Name string
	AttachStdin bool
	AttachStdout bool
	AttachStderr bool
	ExposedPorts nat.PortSet
	Cmd []string
	Image string
	Cpu float64
	Memory int64
	Disk int64
	Env []string
	RestartPolicy string
	//temporary remember to remove
	ContainerId string
}

//Docker struct that will hold the configuration to the Docker client API
type Docker struct {
	Client *client.Client
	Config Config //this field will hold all the configuration of the task 
}

//this is a wrapper around the most common information that is useful for callers
type DockerResult struct {
	Error error
	Action string //action step being undertaken
	ContainerId string
	Result string
}


//pulls images from the imagerepo such as DockerHub
func (d *Docker) Run() DockerResult{
	ctx := context.Background()
	reader, err := d.Client.ImagePull(ctx, d.Config.Image, image.PullOptions{})	
	if err != nil {
		log.Printf("Error pulling image %s: %v\n", d.Config.Image, err)
		return DockerResult{Error: err, Action: "pull", Result: "failed"}
	}

	//copies the info from the io.ReadCloser to the os.Stdout
	//@todo: Interface this to charmCli...later on in the future
	_,err = io.Copy(os.Stdout, reader)
	
	if err != nil {
		log.Fatalf("Error copying logs to the stdout")
	}

	restartPolicy := container.RestartPolicy{
		Name: container.RestartPolicyMode(d.Config.RestartPolicy),
	}

	//specifices the amount of resources required by the host machine to run the container
	resources := container.Resources{
		Memory: d.Config.Memory, 
		NanoCPUs: int64(d.Config.Cpu * math.Pow(10,9)),
	}

	//Specifies the image, env variables and exposedport to run in the container
	cc := container.Config{
		Image: d.Config.Image,
		Tty: false,
		Env: d.Config.Env,
		ExposedPorts: d.Config.ExposedPorts,
	}

	//wrapper of the host configuration
	hc := container.HostConfig{
		RestartPolicy: restartPolicy,
	 	Resources: resources,
		PublishAllPorts: true,
	}

	//create the container with the specified image, and configuration
	log.Println("creating container start operation...")
	resp, err := d.Client.ContainerCreate(ctx, &cc, &hc, nil, nil, d.Config.Name)
	if err != nil {
		log.Printf("Error creating container using image %s: %v\n", resp.ID, err)
		return DockerResult{Error: err, Action: "create", Result: "failed"}
	}

	//start the container
	log.Println("starting container operation...")
	err = d.Client.ContainerStart(ctx, resp.ID, container.StartOptions{})
	if err != nil {
		log.Printf("Error starting container %s: %v\n", resp.ID, err)
		return DockerResult{Error: err}
	}

	//copy the container logs to the stdout
	d.Config.ContainerId = resp.ID 
	

	//gets the container logs dispatched by the container and attaches it to the stdout
	log.Println("switching container logs to the os.Stdout & os.Stderr")
	out, err := d.Client.ContainerLogs(ctx, resp.ID, container.LogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		log.Printf("Error getting logs for container %s: %v\n", resp.ID, err)
		return DockerResult{Error: err}
	}
	
	multiWriter := io.MultiWriter(os.Stdout, os.Stderr)	
	//copy the contents of reader to both output
	_, err = io.Copy(multiWriter, out)
	if err != nil {
		log.Printf("Error copying logs to multiple outpts: %v", err)
		return DockerResult{Error: err}
	}

	return DockerResult{
		ContainerId: resp.ID,
		Action: "start",
		Result: "success",
		Error: nil,
	}
}
 
func (d *Docker) Stop(containerId string) DockerResult {
	log.Printf("Attempting to stop container %v", containerId)
	ctx := context.Background()
	
	//stopping the main container Background Process
	err := d.Client.ContainerStop(ctx, containerId, container.StopOptions{})
	if err != nil {
		log.Printf("Error stopping container %s: %v\n", containerId, err)
		return DockerResult{Error: err, Action: "stopping", Result: "failed"}
	}

	//here we are removing the container main background process
	err = d.Client.ContainerRemove(ctx, containerId, container.RemoveOptions{
		RemoveVolumes: true,
		RemoveLinks: false,
		Force: false,
	})
	if err != nil {
		log.Printf("Error removing container %s: %v\n", containerId, err)
		return DockerResult{Error: err, Action: "removing", Result: "failed"}
	}

	return DockerResult{Error: err, Action: "stop", Result: "success"}


}

func NewConfig(t *Task) *Config {
	//for now we are just dummy typing the config....later on support ui/ux better
	return &Config{
		Name: t.Name,	
		Image: t.Image,
		RestartPolicy: t.RestartPolicy,
		ExposedPorts: t.ExposedPorts,
		//@todo: counter check if the containerId is the same as the Image Id
	}
}

func NewDocker(cfg Config) (*Docker, error) {
	//generate a new client docker request
	dockerClient, err:= client.NewClientWithOpts()
	if err != nil {
		log.Printf("error while creating docker client %v", err)
		return nil, err
	}
	err = client.FromEnv(dockerClient)
	if err != nil {
		log.Printf("error while setting env for client conn %v", err)
		return nil, err
	}

	return &Docker{
		Client: dockerClient,
		Config: cfg,
	}, nil
} 
