package server

import (
	"log"
	"time"

	uuid "github.com/google/uuid"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
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
type Conn struct {
	Id string
	//@todo swap this to a string slice of room ids
	Rooms map[string]*Room //will keep track of the engage rooms by the connection 
	activeRoom *Room //will keep track of which room the conn is currently active
	status status
	Conn *websocket.Conn
	Db *mongo.Collection
	send chan *Message
}

type ClientDto struct {
	Id bson.ObjectID `bson:"_id"`
	Rooms map[string]*Room `bson:"connected_room"`
	ActiveRoom *Room `bson:"active_room"`
	Status string `bson:"status"`
}

func NewConn(conn *websocket.Conn) (*Conn){
	id := uuid.NewString()

	log.Println("creating a new connection...")

	connection := &Conn{
		Id: id,
		Rooms: make(map[string]*Room),
		Conn: conn,
		Db: nil,
		activeRoom: nil,
		status: Online,
		send: make(chan *Message),
	}
	

	log.Println("connection succesfull established. Have a nice chat!")
	return connection
}

func (client *Conn) ConnectDb(db *mongo.Database){
	log.Println("connection to database")
	client.Db = db.Collection("clients")
}
func (client *Conn) DisconnectDb(){
	if client.Db != nil{
		log.Println("disconnecting connection")
		client.Db = nil
	}
	log.Println("connection already disconnected")
}


func (client *Conn) AttachToRoom(rn string) (*Room){
	if _,ok := client.Rooms[rn]; !ok{
		log.Println("room does not exist...creating one")
		//now one can create the room 
		nr := NewRoom(rn)
		//register to the room
		log.Println("registering the client to the room")
		nr.conn[client.Id] = client
		client.appendRoom(nr)
		client.setActiveRoom(nr)

		return nr
	}
	
	room := client.Rooms[rn]
	log.Println("registering the client to the roo")
	room.conn[client.Id] = client

	log.Println("setting current room active")
	client.setActiveRoom(room)

	return room
}

func (client *Conn) DetachToRoom(rn string){
	if _,ok := client.Rooms[rn]; !ok{
		log.Println("room does not exit")
		return
	}

	room := client.Rooms[rn]	
	room.unregister <- client
	client.activeRoom = nil
}

func (client *Conn) appendRoom(room *Room){
	log.Println("appending room to the existing client rooms.")
	client.Rooms[room.Id] = room
}

func (client *Conn) setActiveRoom(room *Room){
	log.Println("setting the client active room.")
	client.activeRoom = room
}

func (client *Conn) UpdateConnStatus(s status){
	client.status = s
}

func (client *Conn) WriteOnceConn(msg []byte) (error) {
	defer client.Conn.Close()
	err := client.Conn.WriteMessage(websocket.TextMessage, msg)
	return err
}

func (client *Conn) ReadOnceConn() ([]byte, error) {
	defer client.Conn.Close()
	var msg []byte
	_, reader, err:= client.Conn.NextReader()
	reader.Read(msg)

	if err != nil{
		return nil, err
	}
	
	return msg, nil
}

func (client *Conn) Serialize() (ClientDto){
	dto := ClientDto{}
	return dto
}

func(client *Conn) ReadMessage(){
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
		client.activeRoom.messages.appendMessage(client, message)
		client.activeRoom.broadcast <- newMsg 
	}
}

func (client *Conn) WriteMessage(){
	defer client.Conn.Close()

	if _,ok :=<-client.send; !ok{
		log.Println("client send channel has been closed")
	}
	client.Conn.SetWriteDeadline(time.Now().Add(writeDeadline))

	for message := range client.send {

		log.Println("sending message to client")
		writer, err := client.Conn.NextWriter(websocket.TextMessage)			


		if err != nil{
			//check if the error is an unexpected error
			if websocket.IsUnexpectedCloseError(err, websocket.CloseAbnormalClosure, websocket.CloseGoingAway) {
				log.Printf("unexpeceted error while writing to the connection: %v", err)
				break
			}
			break
		}

		_,err = writer.Write(message.data)

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseAbnormalClosure, websocket.CloseGoingAway){
				log.Printf("unexpected error while writing to the connection: %v", err)
				break
			}	
			break
		}
	}

}



