package main

import (
	"fmt"
	"log"
	"time"

	"github.com/Chris-Mwiti/build-your-own-x/go_projects/orchestra/manager"
	"github.com/Chris-Mwiti/build-your-own-x/go_projects/orchestra/node"
	"github.com/Chris-Mwiti/build-your-own-x/go_projects/orchestra/task"
	"github.com/Chris-Mwiti/build-your-own-x/go_projects/orchestra/worker"
	"github.com/docker/docker/client"
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
)

func createContainer () (*task.Docker, *task.DockerResult) {
	c := task.Config{
		Name: "test-container-w",
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

	taskOne := task.Task{
		ID: uuid.New(),
		Name: "Task-one",
		State: task.Pending,
		Image: "Image One",
		Memory: 1024,
		Disk: 1,
	}

	taskEvent := task.TaskEvent{
		ID: uuid.New(),
		State: task.Pending,
		Task: taskOne,
		Timestamp: time.Now(),
	}

	fmt.Printf("task: %v\n", taskOne)
	fmt.Printf("task event: %v\n", taskEvent)

	worker := worker.Worker{
		Name: "worker-1",
		Queue: *queue.New(),
		Db: make(map[uuid.UUID]*task.Task),
	}

	fmt.Printf("struct worker: %v\n", worker)
	worker.RunTask()
	worker.StartTask()
	worker.StopTask()
	worker.CollectStats()
	
	manager := manager.Manager{
		Pending: *queue.New(),
		TasksDb: make(map[string][]*task.Task),
		TasksEventDb: make(map[string][]*task.TaskEvent),
		Workers: []string{worker.Name},
		WorkerTaskMap: make(map[string][]uuid.UUID),
	}

	fmt.Printf("struct manager: %v\n", manager)
	manager.SelectWorker()
	manager.SendWork()
	manager.UpdateTask()


	node := node.Node{
		Name: "Node-1",
		Ip: "192.168.1.0",
		Cores: 4,
		Memory: 1024,
		Disk: 25,
		Role: "Worker",
	}

	fmt.Printf("struct node: %v\n", node)
}
