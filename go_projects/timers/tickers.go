package main

import (
	"fmt"
	"time"
)

func Tickers(){
	//used primarly when you want to execute smth repeatedly at regular intervals

	//creation of ticker
	ticker := time.NewTicker(500 * time.Millisecond)
	//creation of a notification channel
	done := make(chan bool)

	go func(){
		for {
			select {
			case <- done:
				fmt.Println("job has been completed")
				return
			case val := <- ticker.C:
				fmt.Println("The value of the ticker is: ", val)
			}
		}
	}()

	//wait for a while for the ticker to be executed
	time.Sleep(1400 * time.Millisecond)
	//Stop the execution of the ticker
	ticker.Stop()

	//send notification to the done channel
	done <- true
	
	//basic simulation of a sleep event to notify that an end of ticker event
	time.Sleep(700 * time.Millisecond)
	fmt.Println("Ticker has stopped")

}