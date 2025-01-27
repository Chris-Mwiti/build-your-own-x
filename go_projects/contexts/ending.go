package main

//interface -> context.Context
//QueryContext -> sql.DB

import (
	"context"
	"fmt"
	"time"
)

type WorkResult string

func endingSomething(ctx context.Context) {
	//ending a context with WithCancel method
	//ctx, cancelCtx := context.WithCancel(ctx) 
	
	//Ending a context with a deadline
	deadline := time.Now().Add(3500 * time.Millisecond);
	deadlineCtx, deadlineCancelCtx := context.WithDeadline(ctx, deadline)

	defer deadlineCancelCtx()

	//just experimenting a usecase: when u use a defer keyword statement to call cancelCtx
	//which is the parent ctx of deadlince ctx a deadlock is produced


	//create an int channel that will be used by the goroutine
	printCh := make(chan int)
	//create a new go routine
	go endingAnother(deadlineCtx, printCh)

	for num := 1; num <= 3; num++ {
		select {
		case <- ctx.Done():
			break
		//pipe the num value to the printCh channel
		case printCh <- num:
			//the thread to sleep for 1 second be resuming with execution
			time.Sleep(2 * time.Second)
		}
}

	deadlineCancelCtx()



	fmt.Printf("ending something \n")
}

func endingAnother(ctx context.Context, printCh <- chan int) {
	for {
		select {
		case <- ctx.Done():
			if err := ctx.Err(); err != nil {
				fmt.Printf("endingAnother err: %s\n", err)
			}
			fmt.Printf("endingAnother: finished\n")
			return
		case results := <-printCh:
			fmt.Printf("endingAnother: %d\n", results)
		}
	}
}

//Advantages of context:
//1. ability to access data stored inside a context

func ending() {
	//one of the two ways of creating a context
	//used as a placeholder when you're not sure which context to use
	ctx := context.TODO()

	//second way of creating contexts
	//used to start a known context
	ctx = context.Background()

	//create an busy loop that will continuosly check if a 
	ctx = context.WithValue(ctx, "myKey", "myValue")
	endingSomething(ctx)
}
