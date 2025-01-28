package main

import "fmt"


func ChannelNonBlocking(){

	//creations of the messages channel and signals channel
	messages := make(chan string)
	signals := make(chan string)

	select {
	case message := <- messages:
		fmt.Printf("It can all be yours chris: %s\n", message)
	default:
		fmt.Println("no messages received")
	}

	msg := "Hello this is about creating non blocking operations"
	select {
	case messages <- msg:
		fmt.Printf("Well a message has been sent:%s\n", &messages)
	default: 
		fmt.Println("no messages have been sent")
	}

	select {
	case signal :=<- signals:
		fmt.Println("received signal", signal)
	case msg := <- messages:
		fmt.Println("received message", msg)
	default:
		fmt.Println("no activity")
	}

}