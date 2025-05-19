package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
	"github.com/gorilla/websocket"
)

func RunServer(){
	//create a new mux handler
	muxHandler := http.NewServeMux()

	baseCtx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	upgrader := &websocket.Upgrader{
		WriteBufferSize: 1024,
		ReadBufferSize: 1024,
	}

	muxHandler.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request){
		//upgrade the connection	
		conn, err := upgrader.Upgrade(w,r,r.Header)
		if err != nil {
			log.Panicf("error while upgrading connection: %v", err)
		}

		for {
			msgType, msg, err := conn.ReadMessage()

			if err != nil {
				log.Println("error while reading in the connection")
				break
			}

			//write the received message
			w, err := conn.NextWriter(msgType)
			if err != nil {
				log.Fatalf("error while writing to connection: %v", err)
			}

			fmtMsg := fmt.Sprintf("received the following message: %s",string(msg))
			w.Write([]byte(fmtMsg))
		}

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
	err := server.ListenAndServe()

	if err != nil {
		log.Fatalf("server shutdown: %v", err)
	}
}


