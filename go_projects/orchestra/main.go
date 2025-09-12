package main

import (
	"fmt"
	"log"
	"time"

	"github.com/Chris-Mwiti/build-your-own-x/go_projects/orchestra/task"
	"github.com/Chris-Mwiti/build-your-own-x/go_projects/orchestra/worker"
	"github.com/docker/docker/client"
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
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
	
	workerDb := make(map[uuid.UUID]*task.Task)
	worker := worker.Worker{
		Db: workerDb,
		Queue: *queue.New(),
	}

	//create a demo task
	tsk := task.Task{
		ID: uuid.New(),
		Name: "test-container-2",
		State: task.Scheduled,
		Image: "strm/helloworld-http",
	}

	//it's a silly algorithim...later to be improved in the future
	log.Println("starting task")
	worker.AddTask(tsk)
	result := worker.Run()
	if result.Error != nil {
		panic(result.Error)
	}

	log.Printf("task %s is runnig in container %s\n", tsk.ID, tsk.ContainerId)
	fmt.Println("sleeping mode...")
	time.Sleep(time.Second * 30)

	log.Printf("stopping task %s\n", tsk.ID)
	tsk.State = task.Completed
	worker.AddTask(tsk)
	result = worker.RunTask()
	if result.Error != nil {
		panic(result.Error)
	}
}

