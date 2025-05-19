package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	uuid "github.com/google/uuid"

	"github.com/gorilla/websocket"
)

//keeps track of the clients status
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
//@todo:
//1. Implement the client connections to be registered in a database(mongoDb)
//2. Ability for the server to implement a server reconnection when needed
//3. Rename the send channel to receive coz why not
//4. 
type ClientConn struct {
	Id string
	Rooms map[string]*Room //will keep track of the engage rooms by the connection 
	activeRoom *Room //will keep track of which room the conn is currently active
	status status
	Conn *websocket.Conn
	send chan *Message
}

func establishConnection(url *url.URL) (*websocket.Conn,*context.Context, error){
	//create a connection string
	s := url.String()

	//create a new context that will be will timeoutes and cancellation
	connCtx, cancel := context.WithTimeout(context.Background(), time.Second * 10)

	//connection dialer
	dialer := websocket.Dialer{
		WriteBufferSize: 1024,
		ReadBufferSize: 1024,
	}
	conn,res,err := dialer.DialContext(connCtx, s, nil)

	if err != nil {
		cancel()
		return nil, nil, err
	}

	log.Printf("connection established succesfully: status code: %s", res.Status)

	return conn, &connCtx, nil	
}

func NewConn(room string, w http.ResponseWriter, r *http.Request) (*ClientConn, error){
	id := uuid.NewString()

	log.Println("establishing a new connection...")
	url := &url.URL{
		Scheme: "ws",
		Host: "localhost:8080",
		Path: "/ws",
	}
	conn,_, err := establishConnection(url)

	if err != nil {
		log.Println("error has occured while establishing connection")
		return nil, fmt.Errorf("error while establishing conn: %v",err)
	}
	

	roomConn := &ClientConn{
		Id: id,
		Conn: conn,
		Rooms: make(map[string]*Room),
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
		log.Println("registering the client to the room")
		nr.conn[client] = client.status

		client.appendRoom(nr)
		client.setActiveRoom(nr)

		return nr
	}
	
	room := client.Rooms[rn]
	log.Println("registering the client to the roo")
	room.conn[client] = Online

	log.Println("setting current room active")
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
	client.activeRoom = nil
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
	client.Conn.SetPongHandler(func(appData string) error {client.Conn.SetReadDeadline(time.Now().Add(setPingWait)); return nil})

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
		client.activeRoom.messages.appendMessage(client, message)
		client.activeRoom.broadcast <- newMsg 
	}
}

func (client *ClientConn) WriteMessage(){
	if _,ok :=<-client.send; !ok{
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
		for i := len(client.send); i >= 0; i--{
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


