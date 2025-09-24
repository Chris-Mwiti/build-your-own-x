package manager

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Chris-Mwiti/build-your-own-x/go_projects/orchestra/task"
	"github.com/Chris-Mwiti/build-your-own-x/go_projects/orchestra/worker"
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
)

var ERR_TASK_404 = errors.New("Task not Found")

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
	CurrentWorker string
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

	manager.CurrentWorker = manager.Workers[newWorker]
	return manager.CurrentWorker
}

//actions: keep track of the resourece stats of the workers
func (manager *Manager) UpdateTask(){
 
	for _, w := range manager.Workers {
		url := fmt.Sprintf("http://%s/tasks", w)

		//@todo: Implement a retry func that will retry failed requests for a number of times
		resp, err := http.Get(url)
		if err != nil {
			log.Printf("error while making get request to url %s: %v\n", url, err)
		}

		if resp.StatusCode != http.StatusOK {
			log.Printf("http error: %d\n", resp.StatusCode)
		}

		var tasks []*task.Task

		//for each task update is status
		de := json.NewDecoder(resp.Body)
		err = de.Decode(&tasks)
		if err != nil {
			log.Printf("error while decoding worker response")
		}

		for _, tsk := range tasks{
			if _, ok := manager.TasksDb[tsk.ID]; !ok{
				log.Println("task not found")
				return
			} 
			manager.TasksDb[tsk.ID] = tsk
		}
	}
}

//actions: add tasks to the task queue
func (manager *Manager) SendWork() (error){

	if manager.Pending.Len() > 0 {
		//dequeue the last added item
		item := manager.Pending.Dequeue()

		//coerce the type to support of type *Task.TaskEvent 
		taskEvent, ok := item.(task.TaskEvent)
		if !ok {
			log.Printf("error while coercing type %v\n", !ok)
			return errors.New("unable to coerce type") 
		}

		taskItem := taskEvent.Task
		//select the worker using a round robin scheduling algorithim
		//@todo: later on in the future we are going to improve the algo
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
			var workerResult worker.ErrResponse

			//decode the resp body to the worker result
			err := json.NewDecoder(resp.Body).Decode(&workerResult)
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

func (manager *Manager) StopWork(taskId uuid.UUID) (error){

	//fetch and make a copy of the task from the task db
	tsk, ok := manager.TasksDb[taskId]
	taskCpy := *tsk
	//adjust the state of the task to completed
	taskCpy.State = task.Completed
	if !ok {
		log.Println("task not found")
	}
	//find which worker is responsible for the assigned task
	wrk, ok := manager.TaskWorkerMap[taskId]
	if !ok {
		return ERR_TASK_404
	}

	//set the current worker to point the registered worker
	manager.CurrentWorker = wrk
	url := fmt.Sprintf("http://%s/tasks",wrk)

	taskEvent := task.TaskEvent{
		ID: uuid.New(),
		State: task.Runnig,
		Task: taskCpy,
		Timestamp: time.Now(),
	}
	
	data, err := json.Marshal(taskEvent)
	if err != nil {
		log.Printf("error while marshaling request %v\n", err)
		return errors.New("Marshaling error")
	}

	res, err := http.Post(url, "application/json", bytes.NewBuffer(data))

	if err != nil {
		log.Printf("error while posting request %v\n", err)
		return errors.New("Error while posting request")
	}

	if res.StatusCode != http.StatusCreated {
		var workerRes worker.ErrResponse
		log.Printf("post request failed %v\n", res.StatusCode)
		err := json.NewDecoder(res.Body).Decode(&workerRes)
		if err != nil {
			return fmt.Errorf("error while decoding res body %v\n", err)
		}

		return errors.New(workerRes.Msg) 	
	}

	return nil
}

func (manager *Manager) AddTask(te task.TaskEvent){
	//responsible to endquer items to the manager queue
	manager.Pending.Enqueue(te)
}

//responsible for listening for any updated task state events
//from listed worker in the worker list
func (manager *Manager) ListenToUpdates(){
	log.Printf("Updating the workers tasks %d\n", len(manager.TasksDb))
	for {
		manager.UpdateTask()
		time.Sleep(15 * time.Second)
	}
}

//construction funcion for managers
func New(workers []string) (*Manager){

	//preenter the workers mapping to the task
	workerTaskmap := make(map[string][]uuid.UUID)
	for _, worker := range workers {
		workerTaskmap[worker] = []uuid.UUID{}
	}

	return &Manager{
		TasksDb: make(map[uuid.UUID]*task.Task),
		TasksEventDb: make(map[uuid.UUID]*task.TaskEvent),
		WorkerTaskMap: workerTaskmap,
		TaskWorkerMap: make(map[uuid.UUID]string),
		Pending: *queue.New(),
		Workers: workers,
	}	
}
