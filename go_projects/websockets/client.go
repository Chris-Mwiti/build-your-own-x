package websockets

import (
	"bytes"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	//time allowed to write a message to a peer
	writeWait = 10 * time.Second

	//time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	//send pings to peer with this period. Must be less than pongWait
	pingPeriod = (pongWait * 9) / 10

	//maximum message size allowed from peer
	maxMessageSize = 512
)

var (
	newLine = []byte{'\n'}
	space = []byte{' '}
)

var upgrader = websocket.Upgrader {
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
}

//Client is a middleman between the websocket connection and the hub
type Client struct {
	hub *Hub

	//the websocket connection
	conn *websocket.Conn

	//Buffered channel of outbound messages.
	send chan []byte

}

//readPump pumps messages from the websocket connection to the hub
//application runs readPump in a per connection goroutine. Ensures that there is at most one reader on a connection 
func (client *Client) readPump() {
	defer func(){
		client.hub.unregister <- client
		client.conn.Close()
	}()

	//sets the readmessage limit for that specific connection
	client.conn.SetReadLimit(maxMessageSize)

	//set the deadline for waiting for a read operation
	//once the deadline is met anything that comes after is corrupt & will result to an error
	client.conn.SetReadDeadline(time.Now().Add(pongWait))
	//sets the write deadline for the write operation
	//once the deadline is met anything that comes after is corrupt & will result to an error
	client.conn.SetWriteDeadline(time.Now().Add(writeWait))

	//provides a callback func with the pong data. 
	//in this case we are setting a deadline for the read operation once the readMessage or the nextreader func is triggerd
	client.conn.SetPongHandler(func(appData string) error {client.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil})

	for {
		//used to get the nextReader and from that initiate from that reader
		//and store the message in a buffer
		_, message, err := client.conn.ReadMessage()
		if err != nil {
			//checks if the websocket err is not among the listed websocket err coded
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure){
				log.Printf("error: %v", err)
			}
			break
		}
		//replaces the entire message with spaces instead of newlines
		message = bytes.TrimSpace(bytes.Replace(message, newLine, space, -1))
		//broadcast the entire message to the hub
		client.hub.broadcast <- message
	}
}

//write pump messages from the hub to the websocket connection
func (client *Client) writePump() {
	//create a new ticker that will automatically send messages
	//to the websocket connection in a interval
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		client.conn.Close()
	}()

	for {
		select {
		case message, ok := <- client.send:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				//if the channel is closed send a message to notify the peer the conn is closed
				client.conn.WriteMessage(websocket.CloseMessage, []byte("Hub closed the send channel"))
				return
			}

			w, err := client.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Printf("error while initializing next writer\n")
				return
			}
			_, err = w.Write(message)

			if err != nil {
				log.Printf("error while writing message\n")
				return
			}

			//add queued chat messsages to the current websocket message
			mesLen := len(client.send)
			for i := 0; i < mesLen; i++ {
				w.Write(newLine)
				w.Write(<-client.send)
			}

		case <-ticker.C:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.conn.WriteMessage(websocket.PingMessage, []byte("ping!!")); err != nil{
				return
			}
		}
	}
}