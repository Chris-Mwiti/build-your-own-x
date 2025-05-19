package main

import (
	"time"
	"log"
	uuid "github.com/google/uuid"
)

//@todo: implement the following methods
//1. Create a new Message
type Message struct {
	Id string
	timestamp time.Time
	data []byte
}

type messageChannel struct {
	sender *ClientConn
	message *Message
}


//@todo: implement the following methods
//1. Find a Message.
//2. Delete a Message
//3. Update a meesage
type MessageHub struct {
	hub map[*ClientConn][]*Message
}

func (msgHub *MessageHub) appendMessage(conn *ClientConn, msg []byte){
	id := uuid.New().String()
	msgHub.hub[conn] = append(msgHub.hub[conn], &Message{
		Id: id,	
		timestamp: time.Now(),
		data: msg,
	})
}

func (msgHub *MessageHub) findMessages(conn *ClientConn)([]*Message){
	messages, ok := msgHub.hub[conn]	

	if !ok {
		log.Println("conn not found in the hub")
		return nil
	}
	
	return messages
}

func newChanMessage(client *ClientConn, message []byte) *messageChannel{
	id := uuid.NewString() 
	return &messageChannel{
		sender: client,
		message: &Message{
			Id: id,
			timestamp: time.Now(),
			data: message,	
		},
	}
}


