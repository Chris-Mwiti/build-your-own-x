package server

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

//@todo: implement the following methods
//1. Creation of a new room
//2. Reconnection of a connection from a client connection
//3. Deleting an existing connection
//4. Termination of a room instance
//5. Display the status of the user

//...[Model Instance]...
//1. Connect the room to a database [Creation of a new room, Update the status of a room,
// termination of room instance, reconnection with a room]
//2. Integration with a pub sub model. This allows when a client send data to a conn, the end
//user can be notified on the message sent
//3. Ability to create private chatting rooms which are encrypted.
//4. Ability to create communities with a limit on the number of users
type Room struct {
	Id  primitive.ObjectID
	RoomId string
	Name string 
	Description string
	MaxConn int
	IsPrivate bool
	CreatedAt time.Time
	db *mongo.Collection
	Conns map[string]*Conn
	messages *MessageHub
	broadcast chan *messageChannel
	register chan *Conn
	unregister chan *Conn
}

type RoomDto struct {
	Id primitive.ObjectID `bson:"_id"`
	Name string `bson:"room_name"`
	RoomId string `bson:"room_id"`
	Description string `bson:"room_description"`
	MaxConn int `bson:"room_maxconn"`
	IsPrivate bool `bson:"room_isprivate"`
	CreatedAt time.Time `bson:"room_createdat"`
	Conns interface{} `bson:"room_conns"`
}

//receives the room name as parameter
//used to create a new room of client conn
func NewRoom(rn string,description string,maxconn int, isprivate bool) *Room{ 
	room := &Room{
		Id: primitive.NewObjectID(),
		Name: rn,
		Description: description,
		MaxConn: maxconn,
		IsPrivate: isprivate,
		CreatedAt: time.Now(),
		db: nil,
		Conns: make(map[string]*Conn),	
		messages: &MessageHub{
			hub: make(map[string][]*Message),
		},
		broadcast: make(chan *messageChannel),
		register: make(chan *Conn),
		unregister: make(chan *Conn),
	} 
	return room
}

func (room *Room) ConnectDb(db *mongo.Database){
	log.Println("connection room to database")
	room.db = db.Collection("rooms")
}

func (room *Room) DisconnectDb() {
	room.db = nil
}

func (room *Room) Serialize() (RoomDto){
	conns := make(map[string]Conn)

	//@todo: refactor and modify the code below to reduce runtime
	for id, conn := range room.Conns{
		conns[id] = *conn
	}

	dto := RoomDto{
		Id: room.Id,
		Name: room.Name,
		Conns: conns,
		Description: room.Description,
		MaxConn: room.MaxConn,
		IsPrivate: room.IsPrivate,
		CreatedAt: room.CreatedAt,
	}

	return dto
}

func DeserializeRoom(room *Room, dto RoomDto){
	room.RoomId = dto.RoomId
	room.Id = dto.Id
	room.Description = dto.Description
}

//always listens for incoming messages
func (room *Room) Listen(ctx context.Context)(error){
	defer room.Close()
	for {
		select{
		case rcvMessage, ok:= <-room.broadcast:
			log.Println("broadcasting message to users")
			//utility check
			clientCount := len(room.Conns)

			if !ok {
				log.Println("error while receiving message")
				continue
			}
			//update the status of the broadcaster
			rcvMessage.sender.UpdateConnStatus(Typing)

			//broadcast to the room users someone is typing
			for _,client:= range room.Conns{
				err := client.Conn.SetWriteDeadline(time.Now().Add(writeDeadline))
				log.Printf("client conn: %d", clientCount)
				if err != nil {
					log.Printf("[RoomListen]:error while setting write deadline: %v", err)
					return err
				}
				//send the message to each of the clients
				client.send <- rcvMessage.message
			}

			rcvMessage.sender.UpdateConnStatus(Online)

		//store the newly created user and update the status
		case client,ok := <-room.register:
			if !ok {
				log.Println("error while registering new client")	
				continue
			}
			room.Conns[client.ClientId] = client
			client.UpdateConnStatus(Online)

		//unregister event listener
		case client, ok := <-room.unregister:
			if !ok {
				log.Println("error while unregistering client")
				continue
			}
			close(client.send)
			delete(room.Conns, client.ClientId)

		case <-ctx.Done():
			return nil

		}

	}
}

func (room *Room) Close(){
	defer func(){
	 close(room.register)
	 close(room.unregister)
	 close(room.broadcast)
	}()
	//@todo: add fields to the room to support the status discovery
}


//room model operations

func(room *Room) CreateRoom(orgCtx context.Context)(*mongo.InsertOneResult, error){
	ctx, cancel := context.WithTimeout(orgCtx, 3 * time.Second)	
	defer cancel()

	result, err := room.db.InsertOne(ctx, room.Serialize())
	if err != nil {
		log.Printf("[CreateRoom]: error while creating room: %v",room.Serialize())
		return nil, err
	}
	return result, nil
}

func(room *Room) FindRoom(orgCtx context.Context, filter bson.D)(*Room, error){
	ctx, cancel := context.WithTimeout(orgCtx, 3 * time.Second)
	defer cancel()
	var result Room

	err := room.db.FindOne(ctx, filter).Decode(&result)
	if err != nil{
		log.Printf("[FindRoom]: error while finding room of filter: %v", filter)
		return nil, err
	}

	return &result, err
}

func (room *Room) UpdateRoom(orgCtx context.Context, filter bson.D, update bson.D)(*mongo.UpdateResult, error){
	ctx, cancel := context.WithTimeout(orgCtx, 3 * time.Second)	
	defer cancel()

	result, err := room.db.UpdateOne(ctx, filter, update)

	if err != nil {
		log.Printf("[UpdateRoom]: error while updating room of filter: %v", filter)
		return nil, err
	}

	return result, nil
}

func (room *Room) DeleteRoom(orgCtx context.Context, filter bson.D)(*mongo.DeleteResult,error){
	ctx, cancel := context.WithTimeout(orgCtx, 3 * time.Second)
	defer cancel()

	result, err := room.db.DeleteOne(ctx, filter)
	if err != nil {
		log.Printf("[DeleteRoom]: error while deleting room of filter: %v", filter)
		return nil, err
	} 
	return result, nil
}




