package worker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/Chris-Mwiti/build-your-own-x/go_projects/orchestra/task"
	"github.com/go-chi/chi/v5"
)

const (
	TASK_KEY = "TASK_KEY"
)

type ErrResponse struct {
	Msg string
	Status uint
}

//this is an inbuilt middleware that is able to fetch a task on pre-request,
//and set it up to the request context
func (api *WorkerApi) TaskCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		taskId := chi.URLParam(r, "taskId") 
		task, err := api.Worker.FetchTaskDb(taskId)
		if err != nil {
			if errors.Is(err, TASK_404){
				log.Printf("task not found error %v\n", err)	
				http.Error(w, "Task Not Found", http.StatusNotFound)
				return
			} else {
				log.Printf("internal server errro while fetching task %v\n", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}
		ctx := context.WithValue(r.Context(), TASK_KEY, task)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

//fetches the worker stats
func (api *WorkerApi) GetStats(w http.ResponseWriter, r *http.Request){
	log.Println("received a get stat request")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(api.Worker.Stats)
	if err != nil {
		log.Printf("error while enconding struct %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (api *WorkerApi) CreateTaskApi(w http.ResponseWriter, r *http.Request){
	log.Println("received a create task request")

	decoder := json.NewDecoder(r.Body)

	//here we are disallowing unknow fields from being accepted or decoded
	decoder.DisallowUnknownFields()
	te := task.TaskEvent{}

	err := decoder.Decode(&te)
	if err != nil {
		errMsg := fmt.Sprintf("error while decoding body %s", err.Error())
		log.Println(errMsg)
	  w.WriteHeader(http.StatusBadRequest)	
		errRes := ErrResponse{
			Msg: errMsg,
			Status: http.StatusBadRequest,
		}

		//encode the res and send back to the user
		err = json.NewEncoder(w).Encode(errRes)
		if err != nil {
			http.Error(w, "Internal Server Error: ", http.StatusInternalServerError)
			return
		}
	}

	//add the task to the worker queue for processing
	api.Worker.AddTask(te.Task)
	log.Printf("Task added %v", te.Task.ID)
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(te.Task)
	if err != nil {
		http.Error(w, "Internal Server Error: ", http.StatusInternalServerError)
		return
	}
}

func (api *WorkerApi) GetTasks(w http.ResponseWriter, r *http.Request){
	log.Println("received a fetch all tasks request")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	tasks, err := api.Worker.FetchTasks()
	if err != nil {
		log.Printf("error while fetching tasks %v", err.Error())
		res := ErrResponse{
			Msg: fmt.Sprint("error while fetching task"),
			Status: http.StatusInternalServerError,
		}

		err = json.NewEncoder(w).Encode(res)
		if err != nil {
			http.Error(w, "Internal Server Error: ", http.StatusInternalServerError)
			return
		}	
	}

	log.Println("fetched tasks from the task database")
	err = json.NewEncoder(w).Encode(tasks)
	if err != nil {
		http.Error(w, "Internal Server Error: ", http.StatusInternalServerError)
		return
	}	
}
func (api *WorkerApi) GetTaskByIdApi(w http.ResponseWriter, r *http.Request){
	log.Println("received a get task request")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	ctx := r.Context()

	task, ok := ctx.Value(TASK_KEY).(*task.Task)
	if !ok {
		http.Error(w, "error while coercing task type", http.StatusInternalServerError)
		return
	}

	log.Printf("found task is of image: %s", task.Image)
	err := json.NewEncoder(w).Encode(task)
	if err != nil {
		http.Error(w, "Internal Server Error: ", http.StatusInternalServerError)
		return
	}	
}

func (api *WorkerApi) PutTaskApi(w http.ResponseWriter, r *http.Request){
	taskId := chi.URLParam(r, "taskId")
	log.Printf("received a put task request %s", taskId)
}

func (api *WorkerApi) StopTaskApi(w http.ResponseWriter, r *http.Request){
	ctx := r.Context()
	retriTask, ok := ctx.Value(TASK_KEY).(*task.Task)
	if !ok {
		http.Error(w, "error while coercing task type", http.StatusInternalServerError)
		return
	}
		//create a copy of the task to add it to the processing task
	taskCopy := *retriTask
	taskCopy.State =	task.Completed 
	api.Worker.AddTask(taskCopy)

	log.Printf("task %v has been added to the queue for processing\n", retriTask.ID)
	w.WriteHeader(http.StatusCreated)
	_, err := w.Write([]byte("task has been added to the queue"))
	if err != nil{
		log.Printf("error while writing back request %v", err)
	}
}

//here we are going to setup the entire path matching for the worker tasks path
func (api *WorkerApi) initRouter(){
	api.Router = router()
	api.Router.Route("/tasks", func(r chi.Router) {
		r.Post("/", api.CreateTaskApi)
		r.Get("/", api.GetTasks)
		r.Route("/{taskId}", func(r chi.Router) {
			r.Use(api.TaskCtx)
			r.Get("/", api.GetTaskByIdApi)
			r.Delete("/", api.StopTaskApi)
		})
	})

	//worker stats route
	api.Router.Route("/stats", func(r chi.Router) {
		r.Get("/", api.GetStats)
	})
	
}

//Start the worker api server
func (api *WorkerApi) Start(){
	log.Println("Starting worker server")
	api.initRouter()
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", api.Address, api.Port), api.Router)
	if err != nil {
		log.Panicf("error while starting worker server")
	}
}
