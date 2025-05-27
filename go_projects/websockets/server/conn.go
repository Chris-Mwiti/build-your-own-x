package server

import (
	"context"
	"fmt"
	"log"
	"time"

	uuid "github.com/google/uuid"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	Id primitive.ObjectID
	ClientId string
	Rooms map[string]*Room //will keep track of the engage rooms by the connection 
	activeRoom *Room //will keep track of which room the conn is currently active
	status status
	Conn *websocket.Conn
	Db *mongo.Collection
	send chan *Message
}

type ClientDto struct {
	Id primitive.ObjectID `bson:"_id"`
	ClientId string `bson:"client_id"`
	Rooms map[string]RoomDto `bson:"connected_room"`
	ActiveRoom RoomDto `bson:"active_room"`
	Status string `bson:"status"`
}

func NewConn(conn *websocket.Conn) (*Conn){
	log.Println("creating a new connection...")

	connection := &Conn{
		Id: primitive.NewObjectID(),
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
		nr.conn[client.ClientId] = client
		client.appendRoom(nr)
		client.setActiveRoom(nr)

		return nr
	}
	
	room := client.Rooms[rn]
	log.Println("registering the client to the room")
	room.conn[client.ClientId] = client

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
	client.Rooms[room.RoomId] = room
}

func (client *Conn) setActiveRoom(room *Room){
	log.Println("setting the client active room.")
	client.activeRoom = room
}

func (client *Conn) UpdateConnStatus(s status){
	client.status = s
}

func (client *Conn) WriteOnceConn(msg []byte) (error) {
	err := client.Conn.WriteMessage(websocket.TextMessage, msg)
	return err
}

func (client *Conn) ReadOnceConn() ([]byte, error) {
	var rn []byte
	
	for {
		_,msg,err := client.Conn.ReadMessage()	

		if err != nil{
			log.Println("could not receive the room name")
			break
		}

		if string(msg) != "" {
			rn = msg
			break
		}	

	}

	return rn, nil
}

func (client *Conn) generateId(){
	id := uuid.NewString()
	client.ClientId = id
}

func (client *Conn) Serialize() (ClientDto){
	//deserialze the client rooms into a format to be supported
	serRoooms := make(map[string]RoomDto)
	for id, room := range client.Rooms {
		serRoooms[id] = room.Serialize()
	}

	//generate a new id for the client
	client.generateId()

	dto := ClientDto{
		Id: client.Id,
		ClientId: client.ClientId,
		Rooms: serRoooms,
		ActiveRoom: RoomDto{},
		Status: string(client.status), 
	}
	return dto
}

func(client *Conn) ReadMessage(ctx context.Context)(error){
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
		select{
		case <-ctx.Done():
			return nil	

		default: 
			_,message,err := client.Conn.ReadMessage()

			if err != nil{
				if websocket.IsUnexpectedCloseError(err, websocket.CloseAbnormalClosure, websocket.CloseGoingAway) {
					log.Printf("unexpected error while closing connection: %v",err)
					return err
				}

				log.Printf("closing connection...: %v", err)
				return err
			}

			//construct the messageSendOb
			newMsg := newChanMessage(client,message)
			client.activeRoom.messages.appendMessage(client, message)
			client.activeRoom.broadcast <- newMsg 

		}
	}
}

func (client *Conn) WriteMessage(ctx context.Context)(error){
	defer client.Conn.Close()

	client.Conn.SetWriteDeadline(time.Now().Add(writeDeadline))

	for {
		select{
		case <-ctx.Done():
			return nil
		default:
			for message := range client.send {

				log.Println("sending message to client")
				err := client.Conn.WriteMessage(websocket.TextMessage, message.data)

				if err != nil {
					if websocket.IsUnexpectedCloseError(err, websocket.CloseAbnormalClosure, websocket.CloseGoingAway){
						log.Printf("unexpected error while writing to the connection: %v", err)
						return err
					}	
					log.Printf("error while writing to the connection: %v",err)
					return err
				}
			}
		}
	}	
}

func (client *Conn) Close(ctx context.Context)(error){
	closeCtx, cancel := context.WithCancel(ctx)
	defer func(){
		close(client.send)
		cancel()
		err := client.Conn.Close()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseAbnormalClosure, websocket.CloseGoingAway){
				log.Panicf("[ClientClose]: unexpected error while closing conn: %v", err)
			}
			log.Println("[ClientClose]: error while closing the client connection")
		}
	}()

	filter := bson.D{bson.E{Key: "client_id", Value: client.ClientId}}
	update := bson.D{bson.E{Key: "$set", Value: bson.E{Key: "status", Value: string(Offline)}}}
	_,err := client.UpdateClient(closeCtx, filter, update)	
	if err != nil {
		log.Println("[ClientClose]: error while updating the status of the client")
		return err
	}
	return nil
}

//database operations...
func (client *Conn) CreateClient(orgCtx context.Context)(*mongo.InsertOneResult, error){
	log.Println("creating new client connection....")
	ctx,cancel := context.WithTimeout(orgCtx, time.Second * 3)
	defer cancel()

	result, err := client.Db.InsertOne(ctx, client.Serialize())
	if err != nil {
		log.Printf("[createClient]: encountered error while creating client: \n")
		return nil, err
	}

	log.Printf("[createclient]: successfully inserted result: %v\n", result)

	return result, nil

}

func (client *Conn) FindClient(orgCtx context.Context,filter bson.D)(*ClientDto,error){
	ctx, cancel := context.WithTimeout(orgCtx, time.Second * 3)
	defer cancel()

	var result ClientDto
	err := client.Db.FindOne(ctx, filter).Decode(&result)

	if err != nil {
		log.Printf("[findclient]: encountered error while finding client filter: %v; error: %v\n", filter, err)
		return nil, err
	}

	return &result, nil
}

func (client *Conn) UpdateClient(orgCtx context.Context, filter bson.D, update bson.D)(*mongo.UpdateResult,error){
	ctx, cancel := context.WithTimeout(orgCtx, time.Second * 3)
	defer cancel()

	result, err := client.Db.UpdateOne(ctx, filter,update) 
	if err != nil {
		log.Printf("[updateclient]; error while updating document filter: %v;error: %v;", filter, err)
		return nil, err
	}

	return result, nil
} 

func (client *Conn) DeleteClient(orgCtx context.Context, filter bson.D)(*mongo.DeleteResult, error){
	ctx, cancel := context.WithTimeout(orgCtx, time.Second * 3)
	defer cancel()

	result, err := client.Db.DeleteOne(ctx, filter)
	if err != nil {
		log.Println("[deleteclient]; error while deleting document")
		return nil, err
	}

	return result, nil
}



