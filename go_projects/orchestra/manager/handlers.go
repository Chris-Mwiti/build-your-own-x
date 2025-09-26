package manager

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/Chris-Mwiti/build-your-own-x/go_projects/orchestra/task"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

const TASK_KEY = "task"

func (api *ManagerApi) WorkersCtx(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {		
		key := chi.URLParam(r, "taskId")
		taskId, err := uuid.Parse(key)
		if err != nil {
			log.Printf("failed to parse uuid %v\n", err)
			http.Error(w, "Bad format of task id", http.StatusBadRequest)
			return
		}
		tsk,err := api.Manger.GetTask(taskId)
		if err != nil {
			if errors.Is(err, ERR_TASK_404){
				http.Error(w, "Task not found", http.StatusNotFound)
				return
			}
			log.Printf("error while fetching task %v\n", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	
		reqCtx := r.Context()
		ctx := context.WithValue(reqCtx, TASK_KEY, tsk)
		next.ServeHTTP(w,r.WithContext(ctx))	
	})
}

func (api *ManagerApi) CreateTaskEvent(w http.ResponseWriter, r *http.Request){		
	//extract the body of the req
	var taskEvent task.TaskEvent
	decoder := json.NewDecoder(r.Body)

	decoder.DisallowUnknownFields()
	err := decoder.Decode(&taskEvent)

	if err != nil {
		log.Printf("error while decoding req body %v\n", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}

	//invoke the manager api event	
	api.Manger.AddTask(taskEvent)
	err = api.Manger.SendWork()
	if err != nil {
		log.Printf("error while posting task event to worker %s: %v\n", api.Manger.CurrentWorker, err)
		http.Error(w, "Error while posting task event to worker", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(taskEvent.Task)
}
 
func (api *ManagerApi) StopTaskEvent(w http.ResponseWriter, r *http.Request){
	//assert the taskid is uuuid format
	key := chi.URLParam(r, "taskId")
	taskId, err := uuid.Parse(key)
	if err != nil {
		log.Printf("task id parsing error")
		http.Error(w, "Task Id does not support the correct format", http.StatusBadRequest)
	}

	err = api.Manger.StopWork(taskId)
	if err != nil {
		log.Printf("error while exec stop work func: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	err = json.NewEncoder(w).Encode(taskId)
	if err != nil{
		log.Printf("error while encoding data : %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusAccepted)
}

func (api *ManagerApi) GetTaskEvent(w http.ResponseWriter, r *http.Request){
	task := r.Context().Value(TASK_KEY)	
	
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(task)
	if err != nil {
		log.Printf("err while encoding data: %v\n", err)
		http.Error(w,"Internal server error", http.StatusInternalServerError)
	}
}


func (api *ManagerApi) initRouter(){
	api.Router = router()

	api.Router.Route("/manager", func(r chi.Router) {
		r.Post("/", api.CreateTaskEvent)

		r.Route("/{taskId}", func(r chi.Router) {
			r.Use(api.WorkersCtx)
			r.Get("/", api.GetTaskEvent)
			r.Delete("/", api.StopTaskEvent)
		})
	})
}

func (api *ManagerApi) Start(){
	log.Println("Starting the manager api server")
	api.initRouter()
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", api.Address, api.Port), api.Router)
	if err != nil {
		log.Panicf("error while starting manager server: %v", err)
	}
}
