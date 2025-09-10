package worker

import (
	"context"
	"log"
	"net/http"

	"github.com/Chris-Mwiti/build-your-own-x/go_projects/orchestra/task"
	"github.com/go-chi/chi/v5"
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

func createTaskApi(w http.ResponseWriter, r *http.Request){
	log.Println("received a create task request")
}
func getTaskApi(w http.ResponseWriter, r *http.Request){
	log.Println("received a get task request")
	ctx := r.Context()

	task, ok := ctx.Value("task").(*task.Task)
	if !ok {
		http.Error(w, "error while coercing task type", http.StatusInternalServerError)
		return
	}

	log.Printf("found task is of image: %s", task.Image)
}

func putTaskApi(w http.ResponseWriter, r *http.Request){
	taskId := chi.URLParam(r, "taskId")
	log.Printf("received a put task request %s", taskId)
}

func deleteTaskApi(w http.ResponseWriter, r *http.Request){
	taskId := chi.URLParam(r, "taskId")
	log.Printf("receive a delete task request %s", taskId)
}

//here we are going to setup the entire path matching for the worker path
func Run() {
	wr := router();

	wr.Route("/tasks", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello this is the worker api"))
		})
		r.Post("/", createTaskApi)
		//@todo: Implement a search api that will be triggered on the name of the image
		//r.Get("/search", searchTaskApi)

		r.Route("/{taskId}", func(r chi.Router) {
			r.Use(TaskCtx)
			r.Get("/", getTaskApi)
			r.Put("/", putTaskApi)
			r.Delete("/", deleteTaskApi)
		})

	})
}


