package manager

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/golang-collections/collections/queue"
	"github.com/Chris-Mwiti/build-your-own-x/go_projects/orchestra/task"
)

//info: Defines the Manager Component Structure
type Manager struct {
	Pending queue.Queue //role:keeps tracked of submitted tasks
	TasksDb map[string][]*task.Task //role:keep track of tasks
	TasksEventDb map[string][]*task.TaskEvent //role:keep track of track events submitted
	Workers []string //role:store the cluster of workers
	WorkerTaskMap map[string][]uuid.UUID //role:map workers and their assigned tasks
	TaskWorkerMap map[uuid.UUID]string //role:map tasks and workers
}

//actions: Pick the appropriate worker from a pool of workers based on their resource stats
func (manager *Manager) SelectWorker(){
	fmt.Println("I will be printing some stuff ")
}

//actions: keep track of the resourece stats of the workers
func (manager *Manager) UpdateTask(){
	fmt.Println("I will be updating task in a few")
}

//actions: add tasks to the task queue
func (manager *Manager) SendWork(){
	fmt.Println("Able to send work to the task queue")

}