package worker

import (
	"fmt"

	"github.com/Chris-Mwiti/build-your-own-x/go_projects/orchestra/task"
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
func (worker *Worker) StopTask(){
	fmt.Println("Stoping task...")
}

