package main

import (
	"fmt"
	"time"
)

//great mechanism for controlling resource utilization
//maintaining quality of service
func RateLimiter(){
	//creation of request channel and feeding it with data
	request := make(chan int, 5)
	for i := 1; i <= 5; i++ {
		request <- i
	}

	//close the request channel after all data has been sent
	close(request)

	//creation of limiter ticker that can be used to regulator in rate limiting scheme
	limiter := time.Tick(200 * time.Millisecond)

	//loop over ther date in the request channel since it is closed
	for req := range request {
		//block on a receive from the limter channel before serving each request
		//limit ourselves to 1 request every 200 millisecond
		<-limiter
		fmt.Println("request", req, time.Now())
	}

	//creation of a burstlimiter 
	burstyLimiter := make(chan time.Time, 5)

	//buffer the init elem of the burstyLimiter to 3 val of time.Now constraint
	for i := 0; i < 3; i++ {
		burstyLimiter <- time.Now()
	}

	//creation of goroutine that will ticker values to the burstyLimiter buffer
	go func(){
		for t := range time.Tick(200 * time.Millisecond){
			burstyLimiter <- t
		}
	}()

	//creation of buffered bursty requests
	burstyRequests := make(chan int, 5)
	for i := 1; i <= 5; i++ {
		burstyRequests <- i
	}
	close(burstyRequests)

	//loop over the requests while blocking the burstyChannel rateLimiter
	for req := range burstyRequests {
		<- burstyLimiter
		fmt.Println("request", req, time.Now())
	} 

}