package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)



func serverWs(db *mongo.Database, w http.ResponseWriter, r *http.Request) {
	ctx := context.WithoutCancel(context.Background()) 
	
	upgrader := websocket.Upgrader{
		WriteBufferSize: 1024,
		ReadBufferSize:  1024,
	}

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Fatalf("could not be able to establish connection: %v", err)
	}

	//create a new client connection
	newConn := NewConn(conn)
	//establish a new connection with the db
	newConn.ConnectDb(db)

	//create a record of the connection 
	_, err = newConn.CreateClient(ctx)
	if err != nil{
		log.Panicf("error while creating client: %v", err)
	}


	//write a message to the user requesting for room name to establish connection to
	err = newConn.WriteOnceConn([]byte("enter the room name or create one?"))
	if err != nil {
		log.Panicf("error while requesting room name: %v", err)
	}

	msg, err := newConn.ReadOnceConn()
	if err != nil {
		log.Panicf("error while receiving room name: %v", err)
	}

	rn := string(msg)
	//here we have to simulate an input space where the user is allowed to enter the room name
	room := newConn.AttachToRoom(rn)
	room.ConnectDb(db)

	//create new go routines to receive and write data
	go newConn.ReadMessage()
	go room.Listen(ctx)
	go newConn.WriteMessage()
}

func RunServer() {
	//opening database file
	db, err := serveDb()
	//dont leave any connection hanging once server is shutdown
	ctx := context.Background()
	defer func(){
		err := db.Disconnect(ctx)
		if err != nil{
			log.Fatalf("error while disconnecting to the database: %v", err)
		}
	}()
	if err != nil {
		log.Fatalf("error can not access the database: %v", err)
	}
	//setup the database to be used
	appDb := db.Database("tchat-db")


	//create a new mux handler
	muxHandler := http.NewServeMux()
	baseCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	//handlers for the connection
	muxHandler.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serverWs(appDb, w, r)
	})

	//create a new server and run it up
	server := http.Server{
		Addr:    "localhost:8080",
		Handler: muxHandler,
		BaseContext: func(l net.Listener) context.Context {
			return baseCtx
		},
	}

	log.Println("server is up and running")
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("server shutdown: %v", err)
	}

}

func serveDb() (*mongo.Client, error) {
	//fetch the uri 
	log.Println("fetchint the env url...")
	uri, ok := os.LookupEnv("MONGO_DB_URL")
	if !ok {
		log.Panic("database url not found")
	}

	log.Println("connecting to the database...")
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("error while establishing database connection: %v", err)
	}	

	return client, nil	
}
