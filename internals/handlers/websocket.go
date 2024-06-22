package handlers

import (
	"encoding/json"
	"fmt"
	"gateway/internals/tools"
	md "gateway/model"
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

type Message struct {
	Action string `json:"action"`
	Data   string `json:"data"`
}

func HandlerWS(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalf("Failed to set up WebSocket upgrade: %v", err)
	}
	defer ws.Close()

	// Infinite loop to listen for messages from the client
	for {
		var response Message
		_, messageBytes, err := ws.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break // Exit the loop in case of error
		}

		var msg Message
		if err := json.Unmarshal(messageBytes, &msg); err != nil {
			log.Println(err)
			continue
		}

		switch msg.Action {
		case "logout":
			response.Action = "logout"
			if status := logoutRequest(w, msg.Data); status != http.StatusOK {
				response.Data = "error"
			} else {
				response.Data = "OK"
			}

		case "echo":
			response = Message{Action: "reply", Data: msg.Data}
		// Add more actions as needed
		default:
			log.Println("Unknown action:", msg.Action)
		}

		responseBytes, _ := json.Marshal(response)
		if err := ws.WriteMessage(websocket.TextMessage, responseBytes); err != nil {
			log.Println(err)
		}
	}

}

func logoutRequest(w http.ResponseWriter, sessionID string) int {
	var user User
	user.SessionID = sessionID
	bodyData := md.RequestBody{
		Action: "logout",
		Body:   user,
	}

	resp, err := tools.SendRequest(w, bodyData, "POST", URLauth)
	if err != nil {
		fmt.Println("Internal server error: " + err.Error())
		return 0
	}
	defer resp.Body.Close()

	return resp.StatusCode
}
