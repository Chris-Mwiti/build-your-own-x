package worker

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Chris-Mwiti/build-your-own-x/go_projects/orchestra/task"
	"github.com/Chris-Mwiti/build-your-own-x/go_projects/orchestra/worker"
	"github.com/go-chi/chi/v5"
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
)

//this is an inbuilt middleware that is able to fetch a task on pre-request,
//and set it up to the request context
func TaskCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		taskId := chi.URLParam(r, "taskId")
		task, err := FetchTaskDb(taskId)
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
}
func (api *WorkerApi) GetTaskApi(w http.ResponseWriter, r *http.Request){
	log.Println("received a get task request")
	ctx := r.Context()

	task, ok := ctx.Value("task").(*task.Task)
	if !ok {
		http.Error(w, "error while coercing task type", http.StatusInternalServerError)
		return
	}

	log.Printf("found task is of image: %s", task.Image)
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
			r.Use(TaskCtx)
			r.Get("/", workerApi.GetTaskApi)
			r.Put("/", workerApi.PutTaskApi)
			r.Delete("/", workerApi.DeleteTaskApi)
		})

	})
}


