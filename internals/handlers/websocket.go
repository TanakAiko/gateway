package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandlerWS(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalf("Failed to set up WebSocket upgrade: %v", err)
	}
	defer ws.Close()

	// Infinite loop to listen for messages from the client
	for {
		messageType, message, err := ws.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break // Exit the loop in case of error
		}

		// Echo the received message back to the client
		if err := ws.WriteMessage(messageType, message); err != nil {
			log.Printf("Error writing message: %v", err)
			break // Exit the loop in case of error
		}
	}

}
