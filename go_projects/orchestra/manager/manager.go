package manager

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/Chris-Mwiti/build-your-own-x/go_projects/orchestra/task"
	"github.com/Chris-Mwiti/build-your-own-x/go_projects/orchestra/worker"
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
)

//info: Defines the Manager Component Structure
type Manager struct {
	Pending queue.Queue //role:keeps tracked of submitted tasks
	TasksDb map[uuid.UUID]*task.Task //role:keep track of tasks
	TasksEventDb map[uuid.UUID]*task.TaskEvent //role:keep track of task events submitted
	Workers []string //role:store the cluster of workers
	WorkerTaskMap map[string][]uuid.UUID //role:map workers and their assigned tasks
	TaskWorkerMap map[uuid.UUID]string //role:map tasks and workers

	//implementing a naive scheduling algorithim
	LastWorker int
}

//actions: Pick the appropriate worker from a pool of workers based on their resource stats
//naive schedluing algorithim (round-robin feature)
func (manager *Manager) SelectWorker() (string){
	fmt.Println("selecting appropriate worker (round-robin algorithim)")
	var newWorker int

	if manager.LastWorker + 1 < len(manager.Workers){
		newWorker = manager.LastWorker + 1
		manager.LastWorker++
	} else {
		newWorker = 0
		manager.LastWorker = 0
	}

	return manager.Workers[newWorker]
}

//actions: keep track of the resourece stats of the workers
func (manager *Manager) UpdateTask(){
	log.Printf("Updating task...how are you feeling with this new keyboard")
}

//actions: add tasks to the task queue
func (manager *Manager) SendWork() (error){

	if manager.Pending.Len() > 0 {
		//dequeue the last added item
		item := manager.Pending.Dequeue()

		//coerce the type to support of type *Task.TaskEvent 
		taskEvent, ok := item.(task.TaskEvent)
		if !ok {
			log.Printf("error while coercing type %v\n", ok)
			return errors.New("unable to coerce type") 
		}

		taskItem := taskEvent.Task
		//select the worker
		selecteWorker := manager.SelectWorker()

		//administrative operations
		manager.TasksEventDb[taskEvent.ID] = &taskEvent 
		manager.WorkerTaskMap[selecteWorker] = append(manager.WorkerTaskMap[selecteWorker], taskItem.ID)
		manager.TaskWorkerMap[taskItem.ID] = selecteWorker
		 

		//adjust the state of the task
		taskItem.State = task.Scheduled
		manager.TasksDb[taskItem.ID] = &taskItem

		log.Println("marshiling the event request for the worker")
		data, err := json.Marshal(taskEvent)
		if err != nil {
			log.Printf("error while marshing the event request %v\n", err)
			return errors.New("Marshiling error")	
		}

		//create a url link to send the request to	
		url := fmt.Sprintf("http://%s/tasks", selecteWorker)
		
		//@todo: later on in the future remove this 
		log.Printf("debugging; generated url (%s)", url)

		resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
		if err != nil {
			log.Printf("error while posting to worker %s: %v\n", url, err)

			//enqueue the faile task for retrial
			manager.Pending.Enqueue(taskEvent)

			return errors.New("error while posting worker")
		}

		//check the statsu code of the response
		if resp.StatusCode != http.StatusCreated {
			log.Printf("response status code: %v\n", resp.StatusCode)	
			workerResult := worker.ErrResponse{} 

			//decode the resp body to the worker result
			err := json.NewDecoder(resp.Body).Decode(workerResult)
			if err != nil {
				log.Printf("error while decoding the response %v\n", err)
				return errors.New("error while decoding the response")
			}

			return errors.New(workerResult.Msg)
		}
	}
	log.Printf("No task event in the queue")
	return nil
}
