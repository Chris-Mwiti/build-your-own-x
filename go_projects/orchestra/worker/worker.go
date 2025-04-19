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
	Queue queue.Queue
	Db map[uuid.UUID]*task.Task
	TaskCount int
}


//responsible for collecting the worker metrics
func (worker *Worker) CollectStats(){
	fmt.Println("Collectiong the task metric")
}

func (worker *Worker) 

