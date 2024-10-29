package ws

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Upgrading the server to handle websockets
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Main handler to handle web socket
// It accepts an instance of Hub and usual gin context
func ServeWs(hub *Hub, c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}
	//creates a new client
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 1024)}
	//send the client to the register channel in the hub
	hub.register <- client

	//Running two go routines write and read to concurrently keep reading and writing messages
	go client.writePump()
	go client.readPump()
}
