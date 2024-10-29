package ws

import "github.com/gorilla/websocket"

// A client will be linked to a hub, will have a socket connection and a send channel to send messages
type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte //send channel is for server to send messages to clients
}

// A hub consists of map to clients, a broadcast channel, a register channel and unregister channel
type Hub struct {
	clients    map[*Client]bool //client connection : T/F; to show if particular client is connected or not
	broadcast  chan []byte      //A channel to broadcast msgs
	register   chan *Client     //Register channel, to register client
	unregister chan *Client     //channel to unregister clients
}

// Creates a new fresh hub
func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

// This function is used to run an infinite loop to run the hub, the main pool
func (h *Hub) Run() {
	for {
		select { //it is like switch
		case client := <-h.register: //Registering a client to hub
			h.clients[client] = true
		case client := <-h.unregister: //Unregistering a client
			if _, ok := h.clients[client]; ok { //in a map, first value is the key, second value ok, shows if he is connected or not which is T or F
				delete(h.clients, client) //If he is connected, delete that client
				close(client.send)        //closes the send messages channel for that client, to free up resources
			}
		case message := <-h.broadcast: //For broadcasting a message
			for client := range h.clients { //a message is received in the form of a byte slice, now we are attempting to send that message to each client
				select { //ranging over all the clients
				case client.send <- message: //if send channel is open, we directly send the message to that client
				default:
					close(client.send)        //else, it means the client is not active,
					delete(h.clients, client) //close the send channel for him and delete him from the hub
				}
			}
		}
	}
}

//Now a client can read also and write also

func (c *Client) readPump() {
	//Defer is there that when this function exits, the client is unregistered from the hub and the conn is closed for the client
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	// Read messages from WebSocket connection
	//An infinite loop to keep reading messages from the client
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		c.hub.broadcast <- message //sending the message to the broadcast channel on the hub, so that server can read messages sent from the client
	}
}

// Used to send messages to the client
func (c *Client) writePump() {
	for {
		select {
		case message, ok := <-c.send: //waiting for a message on the send channel of the client
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.conn.WriteMessage(websocket.TextMessage, message)
		}
	}
}
