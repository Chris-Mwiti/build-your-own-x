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
	wrkshost := os.Getenv("BABOSA_WORKERS_HOST")
	wrksport, _ := strconv.Atoi(os.Getenv("BABOSA_WORKERS_PORT"))
	
	workerDb := make(map[uuid.UUID]*task.Task)
	wrk := worker.Worker{
		Db: workerDb,
		Queue: *queue.New(),
	}


	wrkApi := worker.WorkerApi{
		Address: wrkshost,
		Port: wrksport,	
		Worker: &wrk,
	}

	//listens for incoming or addition of tasks to the queue
	go wrk.Listen()	
	//freq collecting stats
	go wrk.CollectStats()
	//starts the worker http server
	go wrkApi.Start()



	wrks := []string{fmt.Sprintf("%s:%d", wrkshost, wrksport)} 
	mg := manager.New(wrks)


	go mg.Process()
	go mg.ListenToUpdates()
	
	//create a manager api instance
	mngPort,err := strconv.Atoi(os.Getenv("BABOSA_MANAGER_PORT"))
	if err != nil {
		log.Fatalf("port conversion failed")
	}
	mngHost := os.Getenv("BABOSA_MANAGER_HOST")
	mngApi := manager.ManagerApi{
		Port: mngPort,
		Address: mngHost,
		Manger: mg,
	}

	mngApi.Start()
}

