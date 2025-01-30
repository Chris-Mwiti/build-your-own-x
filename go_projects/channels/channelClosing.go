package main

import "fmt"

//closing indicates that no more values
//will be sent on it
//useful to communicate completion to the channel's receivers

func ChannelClosing(){

	//creation of a jobs channel
	jobs := make(chan int, 5)
	done := make(chan bool)

	//creation of a worker go routine
	go func(){
		//more notify's that all the values have been received
		for{
			j, more := <-jobs
			if more {
				fmt.Println("received job", j)
			} else {
				fmt.Println("received all jobs")
				done <- true
				return
			}
		}
	}()

	//Sending jobs to the channel
	for j:=1; j <= 3; j++{
		jobs <- j 
		fmt.Println("sent job",j)
	}
	close(jobs)
	fmt.Println("sent all jobs")

	//synchronization -> wait for all the worker routines to finish
	<-done


	//the receiving end of a channel maps out the following values
	_, ok := <-jobs
	fmt.Println("received more jobs:", ok)
}