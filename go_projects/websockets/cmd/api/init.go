package api

import (
	"context"
	"github.com/gorilla/websocket"
	"net/url"
	"log"
	"time"
)

func EstablishConnection(url *url.URL) (*websocket.Conn,error){
	//create a connection string
	s := url.String()

	//create a new context that will be will timeoutes and cancellation
	connCtx, cancel := context.WithTimeout(context.Background(), time.Second * 10)
	defer cancel()

	//connection dialer
	dialer := websocket.Dialer{
		WriteBufferSize: 1024,
		ReadBufferSize: 1024,
	}
	conn,res,err := dialer.DialContext(connCtx, s, nil)

	if err != nil {
		cancel()
		return nil, err
	}

	log.Printf("connection established succesfully: status code: %s", res.Status)

	return conn, nil	
}


