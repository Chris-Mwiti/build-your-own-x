package manager

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/Chris-Mwiti/build-your-own-x/go_projects/orchestra/task"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func WorkersCtx(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		
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
	
}

