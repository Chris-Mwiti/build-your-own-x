package server

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
	sender *Conn
	message *Message
}


//@todo: implement the following methods
//1. Find a Message.
//2. Delete a Message
//3. Update a meesage
type MessageHub struct {
	hub map[*Conn][]*Message
}

func (msgHub *MessageHub) appendMessage(conn *Conn, msg []byte){
	id := uuid.New().String()
	msgHub.hub[conn] = append(msgHub.hub[conn], &Message{
		Id: id,	
		timestamp: time.Now(),
		data: msg,
	})
}

func (msgHub *MessageHub) findMessages(conn *Conn)([]*Message){
	messages, ok := msgHub.hub[conn]	

	if !ok {
		log.Println("conn not found in the hub")
		return nil
	}
	
	return messages
}

func newChanMessage(client *Conn, message []byte) *messageChannel{
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


