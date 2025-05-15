package websockets

import (
	"errors"
	"log"
	"net/http"
	"time"
	uuid "github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type status string

const (
	Offline status = "Offline"
	Online status = "Online"
	Typing status = "Typing"
)

const (
	readDeadline = 10 * time.Second 
	writeDeadline = 5 * time.Second 	
	setPongWait =  10 * time.Second 
	setPingWait = setPongWait / 4
	readLimit = 256
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
	broadcast chan *messageChannel
	register chan*ClientConn
	unregister chan*ClientConn
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
//2. Attach a conn to a room
//2. Send data over a connection stream
type ClientConn struct {
	Id string
	Rooms []*Room
	activeRoom *Room
	Conn *websocket.Conn
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
		broadcast: make(chan *messageChannel),
	} 
	return room
}

func (room *Room)SetConnStatus(client *ClientConn,conn status) {
	room.conn[client] = conn
}

func NewConn(room string, w http.ResponseWriter, r *http.Request) (*ClientConn, error){
	id := uuid.NewString()


	//@todo: Implent a cli gui form and option selector available rooms

	//new room declaration
	nr := NewRoom(room)
	//update the status of the client in the room

	//upgrade the current http conn to a websocket conn
	upgrader := websocket.Upgrader{
		ReadBufferSize: 1024,
		WriteBufferSize: 1024,
	}

	conn, err := upgrader.Upgrade(w, r, r.Header)

	if err != nil {
		return nil, errors.Join(errors.New("error while establishing connection"), err)
	}
	

	roomConn := &ClientConn{
		Id: id,
		Conn: conn,
		send: make(chan *Message),
	}

	//append the newly created room
	roomConn.appendRoom(nr)
	//set the createdRoom to be the active one
	roomConn.setActiveRoom(nr)

	//adjust the connection of the connection 
	nr.SetConnStatus(roomConn, Online)

	return roomConn, nil
}

func (client *ClientConn) appendRoom(room *Room){
	client.Rooms = append(client.Rooms, room)
}

func (client *ClientConn) setActiveRoom(room *Room){
	client.activeRoom = room
}

func(client *ClientConn) ReadMessage(){
	//set the defaults such as pingtimeouts, ponttimeouts, and close methods
	defer func(){
		err := client.Conn.Close()
		if err != nil {
			log.Printf("error while closing the client connection")
			log.Panic(err)
		}
	}()
	client.Conn.SetReadLimit(readLimit)
	client.Conn.SetPongHandler(func(appData string) error {client.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil})

	//always read for a message
	for{
		_,message,err := client.Conn.ReadMessage()

		if err != nil{
			if websocket.IsUnexpectedCloseError(err, websocket.CloseAbnormalClosure, websocket.CloseGoingAway) {
				log.Printf("error while read websocket: %v", err)
			}
			break
		}

		//construct the messageSendOb
		newMsg := newChanMessage(client,message)
		client.activeRoom.broadcast <- newMsg 
	}
}

func (client *ClientConn) WriteMessage(){
	if _,ok :=<-client.send; ok{
	  log.Panicf("client send channel has been closed")
	}

	for message := range client.send {
			
	}
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


//always listens for incoming messages
func (room *Room) Listen(){
	for {
		select{
		case structMess, ok:= <-room.broadcast:
		if !ok {
				log.Println("error while receiving message")
			}
		//broadcast to the room users someone is typing
		for client, _ := range room.conn{
				err := client.Conn.SetWriteDeadline(time.Now().Add(writeDeadline))
				if err != nil {
					log.Panicf("error while setting write deadline: %v", err)
				}
				//send the message to each of the clients
				client.send <- structMess.message
			}
		}

	}
}



