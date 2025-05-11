package websockets

type Hub struct {
	register chan *Client 
	unregister chan *Client
	broadcast chan []byte
	clients map[*Client]bool
}
