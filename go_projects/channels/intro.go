package main

import "fmt"

func main(){
	//creation of a channel -> using make(chan val -> type)
	//messages channel
	messages := make(chan string)

	//creation of a goroutine thread to pipe data to the channel
	go func() {
		messages <- "ping"
	}()
	
	//receive the channel data in the main thread
	msg := <- messages
	fmt.Println(msg)


	//Other case scenarios of using channels

	//1. Creation of bufferd channels
	bufferedChannels()

	//2. Channels synchronization
	SyncChannel()

	//3. Channels direction
	ChannelDirect()
	
	//4. Channels Select
	SelectChannel()

	//5. Channels Timeouts
	ChannelTimouts()

	//6. Channel non blocking operations
	ChannelNonBlocking()

	//7. Channel closing operations
	ChannelClosing()

	//8. Ranging over channels data
	ChannelRange()

	//9. Workerpools
	WorkerPools()
}