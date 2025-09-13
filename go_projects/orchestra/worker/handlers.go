package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Chris-Mwiti/build-your-own-x/go_projects/orchestra/task"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
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
			log.Printf("error while fetching task %v\n", err)
			http.Error(w, "Internal Server error",http.StatusInternalServerError)
			return
		}
		ctx := context.WithValue(r.Context(), "task", task)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
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
		log.Printf(errMsg)
	  w.WriteHeader(http.StatusBadRequest)	
		errRes := ErrResponse{
			Msg: errMsg,
			Status: http.StatusBadRequest,
		}

		//encode the res and send back to the user
		json.NewEncoder(w).Encode(errRes)
	}

	//add the task to the worker queue for processing

	api.Worker.AddTask(te.Task)
	log.Printf("Task added %v", te.Task.ID)
	w.WriteHeader(http.StatusOK)
	//@todo:for now we will send back the task created although needs to be updated to be more friendly
	json.NewEncoder(w).Encode(te.Task)
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

		json.NewEncoder(w).Encode(res)
	}

	log.Println("fetched tasks from the task database")
	json.NewEncoder(w).Encode(tasks)
}
func (api *WorkerApi) GetTaskByIdApi(w http.ResponseWriter, r *http.Request){
	log.Println("received a get task request")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	ctx := r.Context()

	task, ok := ctx.Value("task").(*task.Task)
	if !ok {
		http.Error(w, "error while coercing task type", http.StatusInternalServerError)
		return
	}

	log.Printf("found task is of image: %s", task.Image)

	json.NewEncoder(w).Encode(task)
}

func (api *WorkerApi) PutTaskApi(w http.ResponseWriter, r *http.Request){
	taskId := chi.URLParam(r, "taskId")
	log.Printf("received a put task request %s", taskId)
}

func (api *WorkerApi) StopTaskApi(w http.ResponseWriter, r *http.Request){
	taskId := chi.URLParam(r, "taskId")
	log.Printf("receive a delete task request %s", taskId)
	if taskId == "" {
		log.Printf("No taskID passed in request. \n")
		w.WriteHeader(http.StatusBadRequest)
	}

	//parse the taskId from a string to a uuid format for retrival
	tId, _ := uuid.Parse(taskId)
	
	//retrival process of the task from the db
	retriTask, ok := api.Worker.Db[tId]

	if !ok {
		log.Printf("task item not availble %v\n", tId)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("task item not found\n"))
	}

	//create a copy of the task to add it to the processing task
	taskCopy := *retriTask
	taskCopy.State =	task.Completed 
	api.Worker.AddTask(taskCopy)

	log.Printf("task %v has been added to the queue for processing\n", taskId)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("task has been added to the queue"))
}

//here we are going to setup the entire path matching for the worker path
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
	
}

func (api *WorkerApi) Start(){
	api.initRouter()
	http.ListenAndServe(fmt.Sprintf("%s:%d", api.Address, api.Port), api.Router)
}
