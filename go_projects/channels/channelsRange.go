package main

import "fmt"

func ChannelRange(){
	//creation of a buffered queue
	queue := make(chan string, 2)
	//sending some data to the key
	queue <- "Queue data 1"
	queue <- "Queue data 2"
	//nb: always remember to close the channel
	close(queue)
	
	//loop over the buffered queue to capture elem
	for elem := range queue {
		fmt.Println(elem)
	}

}