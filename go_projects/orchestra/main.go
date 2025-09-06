package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Chris-Mwiti/build-your-own-x/go_projects/orchestra/task"
	"github.com/docker/docker/client"
)

func createContainer () (*task.Docker, *task.DockerResult) {
	c := task.Config{
		Name: "postgres-container",
		Image: "postgres",
		Env: []string{
			"POSTGRES_USER=cube",
			"POSTGRES_PASSWORD=secrete",
		},
	}

	dc, _ := client.NewClientWithOpts()
	d := task.Docker{
		Client: dc,
		Config: c,
	}
	result := d.Run()

	if result.Error != nil {
		fmt.Printf("%v\n", result.Error)
		return nil, nil
	}

	fmt.Printf("Container %s is running with config %v\n", result.ContainerId, c)

	return &d, &result
}

func stopContainer(d *task.Docker, containerId string) (*task.DockerResult) {
	result := d.Stop(containerId)
	if result.Error != nil {
		log.Printf("Error while stoping container %s: %v\n", containerId, result.Error)
		return &result
	}

	log.Printf("Container %s, has been succesfully stoped\n", containerId)
	return &result
}




func main(){
	fmt.Println("Upcoming orchestra project")

	fmt.Println("Demo: Creating a container")
	dockerTask, createResult := createContainer()
	if createResult.Error != nil {
		fmt.Printf("%v", createResult.Error)
		os.Exit(1)
	}	
	
	//simulate a sleep time
	time.Sleep(2 * time.Second)
	fmt.Printf("Demo: Stopping container %s\n", createResult.ContainerId)

	deleteResult := stopContainer(dockerTask, createResult.ContainerId)
	if deleteResult.Error != nil {
		log.Fatalf("%v\n", deleteResult.Error)
	}
	log.Printf("%v\n", deleteResult)
}
