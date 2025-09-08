package worker

import (
	"fmt"
	"log"
	"time"

	"github.com/Chris-Mwiti/build-your-own-x/go_projects/orchestra/task"
	taskModule "github.com/Chris-Mwiti/build-your-own-x/go_projects/orchestra/task"
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
)

//represent the attribute and methods available for the worker
type Worker struct {
	Name string
	Queue queue.Queue //to keep track of pending tasks
	Db map[uuid.UUID]*task.Task //store the processed tasks and their states
	TaskCount int //count of the tasks assigned to the worker
}


//responsible for collecting the worker metrics
func (worker *Worker) CollectStats(){
	fmt.Println("Collectiong the task metric")
}

//determine the state of a task & actions: start & stop a task based on their state
func (worker *Worker) RunTask(){
	fmt.Println("Running task...")
}

//actions: start a task
func (worker *Worker) StartTask(){
	fmt.Println("Starting task...")

}

//actions: stop a task
func (worker *Worker) StopTask(task *taskModule.Task)(task.DockerResult){
	log.Printf("Stoping task... %s\n", task.ID)
	taskCfg := taskModule.NewConfig(task) 
	dockerClient, err:= taskModule.NewDocker(*taskCfg)
	if err != nil {
		log.Panicf("Panicing: Error while starting docker client: %v", err)
	}
	
	if dockerClient.Config.ContainerId == ""{
		log.Printf("Cannot execute the stop task since the container id is an empty string")
		return taskModule.DockerResult{
			Action: "stop_task",
			Result: "failed:empty_containerId",
			ContainerId: "",
			Error: nil,
		}
	}

	result := dockerClient.Stop(dockerClient.Config.ContainerId)
	if result.Error != nil {
		log.Printf("error while stoping container %s: %v\n", result.ContainerId, result.Error)
		return result
	}


	//here update the state and finish time of the container status
	task.FinishTime = time.Now()	
	task.State = taskModule.Completed	

	log.Printf("Succesfully stopped task %s\n", task.ID)
	return taskModule.DockerResult{
		Action: "stop_task",
		Result: "success",
		ContainerId: result.ContainerId,
		Error: nil,
	}
}

///session2: Concepts Covered
//1. Worker Component Purpose
//2. Define & Implem Algo for Proc Inc tasks
//3. State machine to Transition tasks Btn State


