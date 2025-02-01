package main

import (
	"fmt"
	"sync"
	"time"
)

func WorkerJob(id int){
	fmt.Printf("Worker %d starting\n", id)
	time.Sleep(time.Second)
	fmt.Printf("Worker %d ending\n", id)
}

func WaitGroup(){
	//creation of a wait group instance
	//used to wait for all goroutine launched to finish
	var wg sync.WaitGroup

	//creation of job instances
	for i := 1; i <= 5; i++ {
		
		wg.Add(1)

		go func(){
			defer wg.Done()
			WorkerJob(i)
		}()

	}

	wg.Wait()
}