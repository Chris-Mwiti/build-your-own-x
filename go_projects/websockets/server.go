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
	Rooms map[string]*Room //will keep track of the engage rooms by the connection 
	activeRoom *Room //will keep track of which room the conn is currently active
	status status
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



func NewConn(room string, w http.ResponseWriter, r *http.Request) (*ClientConn, error){
	id := uuid.NewString()
	//@todo: Implent a cli gui form and option selector available rooms
	//upgrade the current http conn to a websocket conn
	upgrader := websocket.Upgrader{
		ReadBufferSize: 1024,
		WriteBufferSize: 1024,
	}

	log.Println("establishing a new connection...")
	conn, err := upgrader.Upgrade(w, r, r.Header)

	if err != nil {
		log.Println("error has occured while establishing connection")
		return nil, errors.Join(errors.New("error while establishing connection"), err)
	}
	

	roomConn := &ClientConn{
		Id: id,
		Conn: conn,
		status: Online,
		send: make(chan *Message),
	}

	log.Println("connection succesfull established. Have a nice chat!")
	return roomConn, nil
}

func (client *ClientConn) AttachToRoom(rn string) (*Room){
	if _,ok := client.Rooms[rn]; !ok{
		log.Println("room does not exist...creating one")
		//now one can create the room 
		nr := NewRoom(rn)
		//register to the room
		log.Println("registering the client to the roo")
		nr.register <- client

		client.appendRoom(nr)
		client.setActiveRoom(nr)

		return nil
	}
	
	room := client.Rooms[rn]
	log.Println("registering the client to the roo")
	room.register <- client
	client.setActiveRoom(room)

	return room
}

func (client *ClientConn) DetachToRoom(rn string){
	if _,ok := client.Rooms[rn]; !ok{
		log.Println("room does not exit")
		return
	}

	room := client.Rooms[rn]	
	room.unregister <- client
}

func (client *ClientConn) appendRoom(room *Room){
	log.Println("appending room to the existing client rooms.")
	client.Rooms[room.Id] = room
}

func (client *ClientConn) setActiveRoom(room *Room){
	log.Println("setting the client active room.")
	client.activeRoom = room
}

func(client *ClientConn) ReadMessage(){
	//set the defaults such as pingtimeouts, ponttimeouts, and close methods
	defer func(){
		log.Println("closing the client connection")
		err := client.Conn.Close()
		if err != nil {
			log.Panicf("error while closing the client connection: %v", err)
		}
	}()
	client.Conn.SetReadLimit(readLimit)
	client.Conn.SetPongHandler(func(appData string) error {client.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil})

	//always read for a message
	for{
		_,message,err := client.Conn.ReadMessage()

		if err != nil{
			if websocket.IsUnexpectedCloseError(err, websocket.CloseAbnormalClosure, websocket.CloseGoingAway) {
				log.Printf("unexpected error while read websocket: %v", err)
			}
			break
		}

		//construct the messageSendOb
		newMsg := newChanMessage(client,message)
		log.Println("broadcasting room message")
		client.activeRoom.broadcast <- newMsg 
	}
}

func (client *ClientConn) WriteMessage(){
	if _,ok :=<-client.send; ok{
	  log.Panicf("client send channel has been closed")
	}
	client.Conn.SetWriteDeadline(time.Now().Add(writeDeadline))

	for message := range client.send {
		writer, err := client.Conn.NextWriter(websocket.TextMessage)			

		defer writer.Close()

		if err != nil{
			//check if the error is an unexpected error
			if websocket.IsUnexpectedCloseError(err, websocket.CloseAbnormalClosure, websocket.CloseGoingAway) {
				log.Panicf("error while writing message: %v", err)
			}
			break
		}

		_,err = writer.Write(message.data)
		if err != nil {
			log.Println("closing the write connection")
			writer.Close()
			log.Panicf("error while writing to the connection: %v", err)
		}

		//send the queued data
		for i := len(client.send); i <= 0; i--{
			message := <-client.send
			_,err = writer.Write(message.data)
			if err != nil {
				log.Println("closing the write connection")
				writer.Close()
				log.Panicf("error while writing to the connections: %v",err)
			}
		}
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
			case rcvMessage, ok:= <-room.broadcast:
			if !ok {
				log.Println("error while receiving message")
			}
			//update the status of the broadcaster
			room.conn[rcvMessage.sender] = Typing

			//broadcast to the room users someone is typing
			for client,_:= range room.conn{
				err := client.Conn.SetWriteDeadline(time.Now().Add(writeDeadline))
				if err != nil {
					log.Panicf("error while setting write deadline: %v", err)
				}
				//send the message to each of the clients
				client.send <- rcvMessage.message
			}

		//store the newly created user and update the status
		case client,ok := <-room.register:
		if !ok {
		 log.Println("error while registering new client")	
		}
		room.conn[client] = Online

	  //unregister event listener
		case client, ok := <-room.unregister:
		if !ok {
				log.Println("error while unregistering client")
		}
		close(client.send)
		delete(room.conn, client)

	}
	}
}



