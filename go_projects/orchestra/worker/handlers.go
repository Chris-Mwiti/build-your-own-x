package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Chris-Mwiti/build-your-own-x/go_projects/orchestra/task"
	"github.com/go-chi/chi/v5"
	"github.com/golang-collections/collections/queue"
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
			http.Error(w, "could not find the task", http.StatusNotFound)
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

func (api *WorkerApi) DeleteTaskApi(w http.ResponseWriter, r *http.Request){
	taskId := chi.URLParam(r, "taskId")
	log.Printf("receive a delete task request %s", taskId)
	
}

//here we are going to setup the entire path matching for the worker path
func Run() {
	wr := router();
	worker := Worker{
		Name: "worker-1",
		Db: make(map[uuid.UUID]*task.Task),
		Queue: *queue.New(),
		TaskCount: 0,
	} 

	taskEvent := task.TaskEvent{
		ID: uuid.New(),
		State: task.Runnig,	
		Timestamp: time.Now(),
		Task: task.Task{
			ID: uuid.New(),
			Name: "test-container-3",
			State: task.Scheduled,
			Image: "strm/helloworld-http",
		},
	} 

	workerApi := WorkerApi{
		Address: "http://localhost",
		Port: "7112",
		Worker: &worker,
	  Router: wr,	
	}


	workerApi.Router.Route("/tasks", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello this is the worker api"))
		})
		r.Post("/", workerApi.CreateTaskApi)
		//@todo: Implement a search api that will be triggered on the name of the image
		//r.Get("/search", searchTaskApi)

		r.Route("/{taskId}", func(r chi.Router) {
			r.Use(workerApi.TaskCtx)
			r.Get("/", workerApi.GetTaskByIdApi)
			r.Put("/", workerApi.PutTaskApi)
			r.Delete("/", workerApi.DeleteTaskApi)
		})

	})
}


