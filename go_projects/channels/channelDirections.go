package main

import "fmt"

//This func send data to a channel using the (ping chan <- msg)

func ping(pings chan <- string, msg string){
	pings <- msg
}

//This func receives data to a channel using the (pong <- chan string)
func pong(pings <- chan string, pongs chan <- string){
	msg := <- pings
	pongs <- msg
}

func ChannelDirect(){
	pings := make(chan string, 1)
	pongs := make(chan string, 1)
	ping(pings, "Hello world")
	pong(pings, pongs)

	fmt.Println(<-pongs)
}