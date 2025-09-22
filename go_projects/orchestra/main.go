package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/Chris-Mwiti/build-your-own-x/go_projects/orchestra/manager"
	"github.com/Chris-Mwiti/build-your-own-x/go_projects/orchestra/task"
	"github.com/Chris-Mwiti/build-your-own-x/go_projects/orchestra/worker"
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)


func init(){
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error could not fetch .env variables %v\n", err)
	}
}

func main(){
	fmt.Println("Upcoming orchestra project")
	host := os.Getenv("BABOSA_HOST")
	port, _ := strconv.Atoi(os.Getenv("BABOSA_PORT"))
	
	workerDb := make(map[uuid.UUID]*task.Task)
	wrk := worker.Worker{
		Db: workerDb,
		Queue: *queue.New(),
	}


	wrkApi := worker.WorkerApi{
		Address: host,
		Port: port,	
		Worker: &wrk,
	}

	//listens for incoming or addition of tasks to the queue
	go wrk.Listen()	
	//freq collecting stats
	go wrk.CollectStats()
	//starts the worker http server
	go wrkApi.Start()


	wrks := []string{fmt.Sprintf("%s:%d", host, port)} 
	mg := manager.New(wrks)

	for i := 0; i < 3; i++{
		tsk := task.Task{
			ID: uuid.New(),
			State: task.Scheduled,
			Name: fmt.Sprintf("test-container-%d", i),
			Image: "strm/helloworld-http",
		}
		te := task.TaskEvent{
			ID:  uuid.New(),
			State: task.Runnig,
			Task: tsk,
		}

		mg.AddTask(te)
		err := mg.SendWork()
		if err != nil {
			log.Printf("error %v", err)
			return 
		}
	} 

	go mg.ListenToUpdates()

	for { 
		for _, tsk := range mg.TasksDb {
			log.Printf("[Manager] Task %s, State %d\n", tsk.ID, tsk.State)
		}
	}

}

