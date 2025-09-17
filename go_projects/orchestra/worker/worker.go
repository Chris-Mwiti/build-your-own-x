package worker

import (
	"errors"
	"fmt"
	"log"
	"slices"
	"time"

	taskModule "github.com/Chris-Mwiti/build-your-own-x/go_projects/orchestra/task"
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
)

var TASK_404 = errors.New("404_TASK")



//state machine for the task state transition
var stateTransitionMap = map[taskModule.State][]taskModule.State{
	taskModule.Pending: {taskModule.Scheduled},
	taskModule.Scheduled: {taskModule.Scheduled,taskModule.Runnig, taskModule.Failed},
	taskModule.Runnig: {taskModule.Runnig,taskModule.Completed,taskModule.Failed},
	taskModule.Completed: {},
	taskModule.Failed: {},
}

//state machine helper functions (engine)
func contains(states []taskModule.State, state taskModule.State) bool {
	if ok := slices.Contains(states, state); ok{
		return true
	}	
	return false
}

func ValidStateTransition(src taskModule.State, dst taskModule.State) bool {
	return contains(stateTransitionMap[src], dst)
}
//represent the attribute and methods available for the worker
type Worker struct {
	Name string
	Queue queue.Queue //to keep track of pending tasks
	Db map[uuid.UUID]*taskModule.Task //store the processed tasks and their states
	Stats *Stats
	TaskCount int //count of the tasks assigned to the worker
}


//responsible for collecting the worker metrics
func (worker *Worker) CollectStats(){
	for {
		log.Printf("fetching stats for worker %s\n", worker.Name)
		worker.Stats = GetStats()
		worker.Stats.TaskCount = worker.TaskCount
		//simulate a sleep func
		time.Sleep(15 * time.Second)
	}
}

func (worker *Worker) Listen(){
	for {
		if worker.Queue.Len() != 0 {
			result := worker.Run()
			if result.Error != nil {
				log.Printf("error while running task %v", result.Error)
			}
		} else {
			log.Printf("no tasks to process currently\n")
		}
		log.Println("sleeping for 15 seconds")
		time.Sleep(15 * time.Second)
	}

}

//determine the state of a task & actions: start & stop a task based on their state
func (worker *Worker) Run() taskModule.DockerResult{
	//here we deque the first task to be uploaded & process
	t := worker.Queue.Dequeue()
	if t == nil {
		log.Println("No tasks in the queue")
		return taskModule.DockerResult{
			Error: nil,
			Action: "worker:dequeue",
			Result: "nil:dequeue",
		}
	}
	//type assertion from interface to task
	taskQueued, ok := t.(taskModule.Task)
	
	if !ok {
		log.Printf("type coercion failed for queued task")
		return taskModule.DockerResult{
			Error: errors.New("error for type coercion failed for queued task"),
			Action: "coercion",
			Result: "coercion:failed",
		}
	}
	//fetch the earlier version
	taskPersisted := worker.Db[taskQueued.ID]

	if taskPersisted == nil {
		taskPersisted = &taskQueued
		worker.Db[taskQueued.ID] = &taskQueued
	}
	var result taskModule.DockerResult

	//@todo: Implement a proper error handling functionality
	if ValidStateTransition(taskPersisted.State, taskQueued.State){
		switch taskQueued.State {
		case taskModule.Scheduled:
			result = worker.StartTask(&taskQueued)
		case taskModule.Completed:
			result = worker.StopTask(&taskQueued)
		default:
			result.Error = errors.New("Not supported state transition")
		}
	} else {
		err := fmt.Errorf("Invalid state transition from %v, to %v", taskPersisted.State, taskQueued.State)
		result.Error = err
	}

	return result
}



//actions: start a task
func (worker *Worker) StartTask(task *taskModule.Task)(taskModule.DockerResult){
	log.Printf("starting task %s", task.ID)

	//get the config of the task
	taskCfg := taskModule.NewConfig(task)
	dockerClient,err := taskModule.NewDocker(*taskCfg)
	if err != nil {
		log.Panicf("Panicing: Error while starting docker client %v", err)
	}

	task.StartTime = time.Now()
	result := dockerClient.Run()

	if result.Error != nil {
		log.Printf("error while running task %v", result.Error)
		task.State = taskModule.Failed
		worker.Db[task.ID] = task
		return result
	}


	log.Println("Running: Completed running task")
	task.State = taskModule.Runnig
	//setting the taskt runtime Id
	task.ContainerId = result.ContainerId

	worker.Db[task.ID] = task
	return result
}

//actions: stop a task
func (worker *Worker) StopTask(task *taskModule.Task)(taskModule.DockerResult){
	log.Printf("Stoping task... %s\n", task.ID)
	taskCfg := taskModule.NewConfig(task) 
	dockerClient, err:= taskModule.NewDocker(*taskCfg)
	if err != nil {
		log.Panicf("Panicing: Error while starting docker client: %v", err)
	}
	
	if task.ContainerId == ""{
		log.Printf("Cannot execute the stop task since the container id is an empty string")
		return taskModule.DockerResult{
			Action: "stop_task",
			Result: "failed:empty_containerId",
			ContainerId: "",
			Error: nil,
		}
	}

	result := dockerClient.Stop(task.ContainerId)
	if result.Error != nil {
		log.Printf("error while stoping container %s: %v\n", result.ContainerId, result.Error)
		return result
	}


	//here update the state and finish time of the container status
	task.FinishTime = time.Now()	
	task.State = taskModule.Completed	

	//update the workers database to keep track of their states
	worker.Db[task.ID] = task

	log.Printf("Succesfully stopped task %s\n", task.ID)
	return taskModule.DockerResult{
		Action: "stop_task",
		Result: "success",
		ContainerId: result.ContainerId,
		Error: nil,
	}
}

//add the task to the queue for execution
func (w *Worker) AddTask(task taskModule.Task) taskModule.DockerResult{
	w.Queue.Enqueue(task)
	return taskModule.DockerResult{
		Action: "add_task",
		Result: "success",
		Error: nil,
	}
}

//dummy prototype of the fetching event
//@todo:implement an algo that wiil be able to conduct a single item search in a queue
func (w *Worker) FetchTaskDb(taskId string) (*taskModule.Task, error) {
	log.Println("fetching task from the datastore")

	//1. parse the taskId to uuid format
	parseId, err := uuid.Parse(taskId)
	if err != nil{
		log.Println("error while parsing the uuid")
		return nil, fmt.Errorf("error while parsing uuid %v", err)
	}

	if task, ok := w.Db[parseId]; ok{
		return task, nil
	}

	return nil, TASK_404 
}

//here for now we are simpling iterating the through an inmemory task db
//@todo: will implement an inbuilt task database like sqlite an sync online with turso
func (w *Worker) FetchTasks() ([]taskModule.Task, error) {
	log.Println("fetching tasks from the datastore")
	var tasks []taskModule.Task

	for _, task := range w.Db{
		tasks = append(tasks, *task)	
	}
	log.Println(tasks)

	return tasks, nil
}


 

///session2: Concepts Covered
//1. Worker Component Purpose
//2. Define & Implem Algo for Proc Inc tasks
//3. State machine to Transition tasks Btn State


