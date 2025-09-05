package task

import (
	"time"
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
}

//Docker struct that will hold the configuration to the Docker client API
type Docker struct {
	Client *client.Client
	Config Config //this field will hold all the configuration of the task 
}
 
