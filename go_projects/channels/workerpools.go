package main

import (
	"fmt"
	"time"
)

//creation of a worker job
func Worker(id int, jobs <-chan int, results chan <- int){
	for j := range jobs {
		fmt.Println("worker", id, "started job", j)
		//simulate an heavy task
		time.Sleep(5 * time.Second)
		fmt.Println("worker", id, "finished job", j)

		//send results to results channel
		results <- j * 2
	}
}

func WorkerPools(){
	//allocate the numboer of jobs
	const numJobs = 5
	//creation of jobs and resuls channel
	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)
	
	for w := 1; w <= 3; w++ {
		go Worker(w, jobs, results)
	}

	for j := 0; j <= numJobs; j++ {
		jobs <- j
	}
	//close the jobs channel after sending data
	close(jobs)

	for a := 1; a <= numJobs; a++ {
		fmt.Println("The following are the results for ", a, ":", <- results)
	}
}