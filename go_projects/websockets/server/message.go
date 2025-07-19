package server

import (
	"context"
	"errors"
	"log"
	"time"

	uuid "github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

//@todo: implement the following methods
//1. Create a new Message
type Message struct {
	Id string
	timestamp time.Time
	data []byte
}

type MessageDto struct {
	Id primitive.ObjectID `bson:"_id"`
	MessageId string `bson:"message_id"`
	data string `bson:"data"`
	CreatedAt time.Time `bson:"createdAt"`
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
	//the hub should store a map of the map of the conn id with the message sent
	RoomId string
	RcvCh chan *messageChannel 
	hub map[string][]*Message
	coll *mongo.Collection
}

type MessageHubDto struct {
	Id primitive.ObjectID `bson:"_id"` 
	RoomId string `bson:"room_id"`
	Messages []MessageDto `bson:"messages"`
}

func (msgHub *MessageHub) appendMessage(clientId string, msg []byte){
	log.Println("appending message to the hub")
	id := uuid.New().String()
	msgHub.hub[clientId] = append(msgHub.hub[clientId], &Message{
		Id: id,	
		timestamp: time.Now(),
		data: msg,
	})
}

func (msgHub *MessageHub) findMessages(clientId string)([]*Message){
	messages, ok := msgHub.hub[clientId]	

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

func (msg *Message) Serialize() (MessageDto) {

	serMsg := MessageDto{
		Id: primitive.NewObjectID(),
		MessageId: msg.Id,
		CreatedAt: msg.timestamp,
		data: string(msg.data),
	}

	return serMsg
}

func (msgHub *MessageHub) Serialize() (MessageHubDto){
	var serializedMessage []MessageDto

	for _, hub := range msgHub.hub{
		for _, msg := range hub {
			serializedMessage = append(serializedMessage, msg.Serialize())
		}	
	}

	return MessageHubDto{
		Id: primitive.NewObjectID(),
		RoomId: msgHub.RoomId,
		Messages: serializedMessage,
	}
}

func (msgHub *MessageHub) Listen(ctx context.Context) (error) {
	for {
		select{
		case <-ctx.Done():
		log.Printf("closing the messaging listening channel")
		
		case msg := <-msgHub.RcvCh:
		log.Printf("received message from client: %v", msg.sender.ClientId)
		msgHub.appendMessage(msg.sender.ClientId, msg.message.data)
		}
	}
	
}

func (msgHub *MessageHub) CreateRoomHub(ctx context.Context) (*mongo.InsertOneResult, error){
	log.Println("Creating room message hub")
	result, err := msgHub.coll.InsertOne(ctx, msgHub.Serialize())
	if err != nil {
		log.Printf("[CreateRoomHub]: error while creating room hub: %v", err)
		return nil, err
	}

	return result, nil
}	

func (msgHub *MessageHub) InsertMessage(ctx context.Context, msg *messageChannel)(*mongo.UpdateResult, error){
	log.Println("Inserting message into room hub")

	filter := bson.D{bson.E{Key: "room_id", Value: msgHub.RoomId}}
	var hub MessageHubDto
	err := msgHub.coll.FindOne(ctx, filter).Decode(&hub)
	if err != nil && errors.Is(err,mongo.ErrNoDocuments){
		//create a new hub for the room	
		msgHub.CreateRoomHub(ctx)
	} else {
		log.Printf("[InsertMessage]: error while inserting message: %v", err)
		return nil, err	
	} 
	
	update := bson.D{bson.E{Key: "messages", Value: msgHub.Serialize().Messages}}
	result, err := msgHub.coll.UpdateOne(ctx, filter, update)

	return result, nil
}

