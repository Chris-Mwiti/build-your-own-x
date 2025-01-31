package main

import (
	"fmt"
	"time"
)

func Timers(){
	//timers are used to represent a single event in the future
	//you tell the timer how long you want to wait and it provides a channel
	//that will be notified at that time
	timer1 := time.NewTimer(2 * time.Second)

	//blocks the timer's channel C until it sends a value indication that the timer fired
	<- timer1.C
	fmt.Println("Timer 1 has fired")

	//creation of a timer that can be stopped
	timer2 := time.NewTimer(4 * time.Second)
	//creation of go routine
	go func(){
		//block the timer2 channel
		<- timer2.C
		fmt.Println("Timer 2 has fired")
	}()

	//stop the timer 2 before even execution
	stop2 := timer2.Stop();
	if stop2 {
		fmt.Println("Timer 2 has stopped")
	}

	time.Sleep(5 * time.Second)

}