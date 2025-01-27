package main

import (
	"fmt"
	"time"
)

func SelectChannel(){
	//Each channel will receive a value after some amount of time
	//to simulate blocking rpc operations execution
	c1 := make(chan string)
	c2 := make(chan string)

	go func(){
		//Sleep the first goroutine to 1 second
		time.Sleep(1 * time.Second)
		c1 <- "Musa"
	}()

	go func(){
		//Sleep the second goroutine to 2 second
		time.Sleep(2 * time.Second)
		c2 <- "Imbeka"
	}()

	for i := 0; i < 2; i++ {
		select {
		case msg1 := <- c1:
			fmt.Println("received", msg1)
		case msg2 := <- c2:
			fmt.Println("received", msg2)
		}
	}


}