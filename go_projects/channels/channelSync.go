package main

import (
	"time"
	"fmt"
)

//goroutine func that will send a done mes to a channel
//to notify another goroutine that this func work is done
func worker(done chan bool){
	fmt.Printf("working...")
	time.Sleep(time.Second)
	fmt.Printf("done\n")

	//send the done message to the channel 
	//to trigger the completion of a task
	done <- true
}

func SyncChannel(){

	//Creation of a buffered channel
	done := make(chan bool, 1)

	//creation of a goroutine
	go worker(done)

	//sync of the goroutines to ensure that the execution 
	//of the main goroutine is waiting till completion of the worker threads
	//if removed the worker thread would not be even started
	<-done

	fmt.Println("worker goroutines are completed")
}