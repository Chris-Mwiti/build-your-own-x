package main

import "fmt"

func bufferedChannels(){
	//creation of a buffered channel
	messages := make(chan string, 2)

	//sending data to the bufferd channel
	messages <- "buffered" 
	messages <- "channel"

	fmt.Println(<-messages)
	fmt.Println(<-messages)

}