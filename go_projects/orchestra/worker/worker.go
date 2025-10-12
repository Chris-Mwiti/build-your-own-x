package worker

import (
	"context"
	"errors"
	"fmt"
	"log"
	"slices"
	"time"

	taskModule "github.com/Chris-Mwiti/build-your-own-x/go_projects/orchestra/task"
	"github.com/Chris-Mwiti/build-your-own-x/go_projects/orchestra/utils"
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
)

var TASK_404 = errors.New("404_TASK")
var COERCION_ERROR = errors.New("COERCION_ERROR")
var TRANSITION_NOT_SUPPORTED = errors.New("TRANSITION_NOT_SUPPORTED")
var ERR_FUNC_EXEC = errors.New("func execution error")



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

//responsible to listen for any incoming task to the queue
func (worker *Worker) Listen(){
	for {
		if worker.Queue.Len() != 0 {
			result := worker.run()
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
func (worker *Worker) run() taskModule.DockerResult{
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
			Error: COERCION_ERROR,
			Action: "coercion",
			Result: "coercion:failed",
		}
	}
	//fetch the earlier version
	taskPersisted := worker.Db[taskQueued.ID]

	//if the task is not persisted
	if taskPersisted == nil {
		taskPersisted = &taskQueued
		worker.Db[taskQueued.ID] = &taskQueued
	}
	var result taskModule.DockerResult

	//@todo: Implement a proper error handling functionality
	//@todo: Later on in the future we are going to implement support for other transitioins
	if ValidStateTransition(taskPersisted.State, taskQueued.State){
		switch taskQueued.State {
		case taskModule.Scheduled:
			result = worker.StartTask(&taskQueued)
		case taskModule.Completed:
			result = worker.StopTask(&taskQueued)
		case taskModule.Failed:
			result = worker.RestartTask(&taskQueued)
		default:
			result.Error = TRANSITION_NOT_SUPPORTED
		}
	} else {
		err := fmt.Errorf("invalid state transition from %v, to %v", taskPersisted.State, taskQueued.State)
		result.Error = err
	}

	return result
}

func (worker *Worker) RestartTask(task *taskModule.Task)(taskModule.DockerResult){
	log.Printf("restartin task %v\n", task.ID)
	_,err := utils.RetryFn(context.Background(), 3, 5 * time.Second, func(ctx context.Context) (*taskModule.Task, error) {
		tsk := worker.AddTask(*task)
		
		if tsk.Error != nil {
			log.Printf("error while adding task to worker queue")
			return nil, tsk.Error
		}
		result := worker.run()
	
		if result.Error != nil {
			log.Printf("error while re running task %v\n", task.ID)
			return nil, result.Error
		}
		
		return task, nil
	})  

	if err != nil{
		log.Printf("error while retrying to restart task %v\n", err)

		return taskModule.DockerResult{
			Error: err,
			Action: "restart",
			Result: "restart:failed",
		}
	}
	

	return taskModule.DockerResult{
		Error: nil,
		Action: "restart",
		Result: "restart:success",
	}
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
		//@todo: implement the util retry functionality
		enhDockerClient := func (ctx context.Context)(taskModule.DockerResult, error){
			if err := ctx.Err(); err != nil{
				log.Println("wrapper func: context done")
				return taskModule.DockerResult{
					Error: err,
				}, nil
			}	
			return dockerClient.Run(), nil
		}
		result, err = utils.RetryFn(context.Background(),3,(4 * time.Second), enhDockerClient)
		if err != nil {
			log.Printf("Failed to retry: [dockerClient.Run]: %v\n", err)
		}
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

func (w *Worker) ListenUpdateTasks() (error) {
	for {
		log.Println("listening to update tasks events")
		err := w.updateTasks()

		if err != nil {
			log.Printf("error while listening to update tasks events: %v\n", err)
			return ERR_FUNC_EXEC 
		}
		log.Printf("sleeping for 15 seconds")
		time.Sleep(15 * time.Second)
	}
}

func (w *Worker) inspectTask(task *taskModule.Task) taskModule.DockerInspectResult{
	config := taskModule.NewConfig(task)
	dockerClient, err := taskModule.NewDocker(*config)
	if err != nil {
		log.Printf("error while configuring new docker client %v\n", err)
		return taskModule.DockerInspectResult{
			Error: err,
		}
	}

	return dockerClient.Inspect(task.ContainerId)
}

func (w *Worker) updateTasks() (error) {
	for id, task := range w.Db {
		//here we will check if the task state is Running
		//and if so get insights and health check about it
		resp := w.inspectTask(task)
		
		if resp.Error != nil {
			log.Printf("error while inspecting container %v\n", resp.Error)
			return fmt.Errorf("error while inspecting container %v\n", resp.Error)
		}

		//if the container is nil then it means it has already failed
		if resp.Container == nil {
			log.Printf("container for task %s in non-runnig state %s\n", task.ID, task.State)
			//update the state of task in the db to point to failing
			w.Db[id].State = taskModule.Failed
			return nil
		}

		//the container has already existed so the state of the task is failed
		if resp.Container.State.Status == "exited" {
			log.Printf("container for task %s in non-running state %s\n", id, task.State)
			w.Db[id].State = taskModule.Failed
		}

		w.Db[id].PortBindings = resp.Container.NetworkSettings.NetworkSettingsBase.Ports
	}	

	return nil
}




 


