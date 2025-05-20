package server

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"
	"github.com/gorilla/websocket"
)


func serverWs(w http.ResponseWriter, r *http.Request){
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
	rn := "text_room"

	//here we have to simulate an input space where the user is allowed to enter the room name
	room := newConn.AttachToRoom(rn)

	//create new go routines to receive and write data
	go room.Listen()
	go newConn.ReadMessage()
	go newConn.WriteMessage()
}

func RunServer(){
	//create a new mux handler
	muxHandler := http.NewServeMux()

	baseCtx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	//handlers for the connection
	muxHandler.HandleFunc("/ws", serverWs)

	
	//create a new server and run it up
	server := http.Server{
		Addr: "localhost:8080",	
		Handler: muxHandler,
		BaseContext: func(l net.Listener) context.Context {
			return baseCtx
		},
	}

	log.Println("server is up and running")
	err := server.ListenAndServe()

	if err != nil {
		log.Fatalf("server shutdown: %v", err)
	}
	
}


