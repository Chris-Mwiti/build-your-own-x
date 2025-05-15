package websockets

import (
	"log"
	"time"

	"github.com/google/uuid"
)

type status string

const (
	Offline status = "Offline"
	Online status = "online"
)

//@todo: implement the following methods
//1. Creation of a new room
//2. Reconnection of a connection
//3. Deleting a new connection
//2. Updating the status of client conn
type Room struct {
	Id string
	Name string
	conn map[*ClientConn]status
	messages *MessageHub
	receive chan *messageChannel
}

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

//@todo: implement the following methods
//1. Create a new conn
//2. Send data over a connection stream
type ClientConn struct {
	Id string
	Rooms []*Room
	send chan *Message
}

//receives the room name as parameter
//used to create a new room of client conn
func NewRoom(rn string) *Room{ 
	id := uuid.NewString();
	room := &Room{
		Id: id,
		Name: rn,
		conn: make(map[*ClientConn]status),	
		messages: new(MessageHub),
		receive: make(chan *messageChannel),
	} 
	return room
}



func NewConn(room string) *ClientConn{
	id := uuid.NewString()

	conn := &ClientConn{
		Id: id,
		Rooms: make([]*Room, 0),
		send: make(chan *Message),
	}

	return conn
}

//always listens for incoming messages
func (room *Room) Listen(){
	for {
		select{
		case message, ok:= <-room.receive:
		if !ok {
				log.Println("error while receiving message")
			}
		//store the message for the client conn in the room
		room.messages = &MessageHub{}		
		}

	}
}
