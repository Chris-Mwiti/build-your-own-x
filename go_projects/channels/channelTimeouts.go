package main

import (
	"fmt"
	"time"
)

func ChannelTimouts() {
	//creation of a new channel
	c1 := make(chan string, 1)

	//creation of a goroutine that wi\nll send data to the channel
	go func(){
		time.Sleep(2 * time.Second)
		c1 <- "Hello chris this is the first channel data\n"
	}()
	
	select {
	case res := <- c1:
		fmt.Println(res)
	case <-time.After(1 * time.Second):
		fmt.Println("The following is a timeout that has been scheduled: 1")
		//close the channel
	}

	//creation of channel 2
	c2 := make(chan string, 1)
	go func(){
		time.Sleep(2 * time.Second)
		c2 <- "Hello chris this is the second channel data\n"
	}()

	//select the channel that actively receives the data
	select {
	case res2 := <- c2:
		fmt.Println(res2)
	case <-time.After(3 * time.Second):
		//close the channel
		close(c2)
		fmt.Println("The following is a timeout that has been scheduled ")
	}
	
}