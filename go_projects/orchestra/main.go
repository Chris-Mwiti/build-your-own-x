package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/Chris-Mwiti/build-your-own-x/go_projects/orchestra/task"
	"github.com/Chris-Mwiti/build-your-own-x/go_projects/orchestra/worker"
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
)




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

	//starts the worker http server
	wrkApi.Start()
	//listens for incoming or addition of tasks to the queue
	go wrk.Listen()	
}

