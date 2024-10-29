package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/raunak173/go-chat/ws"
)

func main() {
	router := gin.Default()
	//Hub is the palce that handles everything. It is where client registers and de registers and the messages are broadcasted
	hub := ws.NewHub() //Creating a fresh new hub
	go hub.Run()       //creating a go routine to run the hub, so that we can run it concurrently

	router.GET("/ws", func(c *gin.Context) {
		ws.ServeWs(hub, c) //serveWs is the main function that handles the websocket
	})

	log.Println("Server started on http://localhost:8000")
	if err := router.Run(":8000"); err != nil {
		log.Fatal("Server failed:", err)
	}
}
