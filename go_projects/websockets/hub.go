package websockets

type Hub struct {
	register chan *Client 
	unregister chan *Client
	broadcast chan []byte
	clients map[*Client] bool
}

func newHub()*Hub{
	hub := &Hub{
		register: make(chan *Client),
		unregister: make(chan *Client),
		broadcast: make(chan []byte),
		clients: make(map[*Client]bool),
	}
	return hub
}

func (hub *Hub) run(){
	for{
		select{

		case client := <-hub.register:
		//add the connection to the clients
		hub.clients[client] = true

		case client := <-hub.unregister:
		close(client.send)
		delete(hub.clients, client)

		case message := <-hub.broadcast:
		for client := range hub.clients{
			select{
				case client.send <-message:
				default:
				close(client.send)
				delete(hub.clients, client)
			}	
		}

		}
	}
}
