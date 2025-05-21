package server

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/boltdb/bolt"
	"github.com/gorilla/websocket"
)

const (
	databaseFile = "./database/storage.db"
)


func serverWs(db *bolt.DB, w http.ResponseWriter, r *http.Request){
	upgrader := websocket.Upgrader{
		WriteBufferSize: 1024,
		ReadBufferSize: 1024,
	}


	conn, err := upgrader.Upgrade(w,r,nil)

	if err != nil {
		log.Fatalf("could not be able to establish connection: %v", err)
	}

	//create a new client connection
	newConn := NewConn(conn) 
	//establish a new connection with the db
	newConn.ConnectDb(db)

	//write a message to the user requesting for room name to establish connection to
	err = newConn.WriteOnceConn([]byte("enter the room name or create one?"))
	if err != nil{			
		log.Panicf("error while requesting room name: %v", err)
	}

	msg, err := newConn.ReadOnceConn()
	if err != nil {
		log.Panicf("error while receiving room name: %v", err)
	}

	rn := string(msg)
	//here we have to simulate an input space where the user is allowed to enter the room name
	room := newConn.AttachToRoom(rn)

	//create new go routines to receive and write data
	go room.Listen()
	go newConn.ReadMessage()
	go newConn.WriteMessage()
}

func RunServer(){
	//opening database file
	db, err := serveDb(databaseFile)
	//dont leave any connection hanging once server is shutdown
	defer db.Close()

	if err != nil {
		log.Fatalf("error can not access the database: %v", err)
	}

	//create a new mux handler
	muxHandler := http.NewServeMux()
	baseCtx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	//handlers for the connection
	muxHandler.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serverWs(db, w, r)
	})

	
	//create a new server and run it up
	server := http.Server{
		Addr: "localhost:8080",	
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

func serveDb(dbFile string) (*bolt.DB, error){
	//check if the file exists
	if _,err := os.Stat(dbFile); err != nil {
		log.Fatalf("error while setting up database: %v: db file does not exist", err)
	}

	file, err := bolt.Open(dbFile,0600,&bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}

	return file, nil
}


