package server

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
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
	Id  string
	Name string 
	db *mongo.Collection
	conn map[string]*Conn
	messages *MessageHub
	broadcast chan *messageChannel
	register chan *Conn
	unregister chan *Conn
}

type RoomDto struct {
	Id bson.ObjectID `bson:"_id"`
	Name string `bson:"room_name"`
}

//receives the room name as parameter
//used to create a new room of client conn
func NewRoom(rn string) *Room{ 
	id := uuid.NewString();
	room := &Room{
		Id: id,
		Name: rn,
		db: nil,
		conn: make(map[string]*Conn),	
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
	//@todo: complet the serialization method
	dto := RoomDto{
		Name: room.Name,
	}

	return dto
}

//always listens for incoming messages
func (room *Room) Listen(ctx context.Context){
	defer func(){
		close(room.broadcast)
		close(room.register)
		close(room.unregister)
	}()
	for {
		select{
		case rcvMessage, ok:= <-room.broadcast:
			log.Println("broadcasting message to users")
			//utility check
			clientCount := len(room.conn)

			if !ok {
				log.Println("error while receiving message")
				continue
			}
			//update the status of the broadcaster
			rcvMessage.sender.UpdateConnStatus(Typing)

			//broadcast to the room users someone is typing
			for _,client:= range room.conn{
				err := client.Conn.SetWriteDeadline(time.Now().Add(writeDeadline))
				log.Printf("client conn: %d", clientCount)
				if err != nil {
					log.Panicf("error while setting write deadline: %v", err)
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
			room.conn[client.Id] = client
			client.UpdateConnStatus(Online)

		//unregister event listener
		case client, ok := <-room.unregister:
			if !ok {
				log.Println("error while unregistering client")
				continue
			}
			close(client.send)
			delete(room.conn, client.Id)
	
		case <-ctx.Done():
			return

		}

	}
}




