package handlers

import (
	"encoding/json"
	"fmt"
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

type Client struct {
	User md.User
	Conn *websocket.Conn
	Send chan []byte
}

var clients = make(map[int]*Client)
var broadcast = make(chan []byte)

func HandlerWS(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("Failed to set up WebSocket upgrade: %v\n", err)
		return
	}
	defer ws.Close()

	sessionID, err := r.Cookie("sessionID")
	if err != nil {
		fmt.Printf("Failed to get cookie: %v\n", err)
		return
	}

	status, user := getUserData(w, sessionID.Value)
	if status != http.StatusOK {
		fmt.Printf("Failed to get user's data: %v\n", err)
		return
	}

	client := &Client{User: user, Conn: ws, Send: make(chan []byte)}
	clients[client.User.Id] = client

	go writeToClients(client)

	// Infinite loop to listen for messages from the client
	for {
		var response Message
		_, messageBytes, err := ws.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			delete(clients, client.User.Id)
			break // Exit the loop in case of error
		}

		var msg Message
		if err := json.Unmarshal(messageBytes, &msg); err != nil {
			log.Println(err)
			continue
		}

		fmt.Println("msg: ", msg)

		switch msg.Action {
		case "logout":
			response.Action = "logout"
			if status := logoutRequest(w, msg.Data); status != http.StatusOK {
				response.Data = "error"
			} else {
				response.Data = "OK"
			}

		case "messageCreate":
			response.Action = "messageCreate"
			if status := createMessage(w, msg.Data); status != http.StatusOK {
				response.Data = "error"
			} else {
				response.Data = "OK"
				broadcast <- []byte(msg.Data)
			}

		case "messageGets":
			response.Action = "messageGets"
			status, chats := getMessages(w)
			if status != http.StatusOK {
				response.Data = "error"
			} else {
				response.Data = chats
			}

		case "messageStatusReceived":
			response.Action = "messageStatusReceived"
			if status := setStatusReceived(w, msg.Data); status != http.StatusOK {
				response.Data = "error"
			} else {
				response.Data = "OK"
			}

		case "messageStatusRead":
			response.Action = "messageStatusRead"
			if status := setStatusRead(w, msg.Data); status != http.StatusOK {
				response.Data = "error"
			} else {
				response.Data = "OK"
			}

		case "postCreate":
			response.Action = "postCreate"
			if status := createPost(w, msg.Data); status != http.StatusCreated {
				response.Data = "error"
			} else {
				response.Data = "OK"
			}

		case "getOnePost":
			response.Action = "getOnePost"
			status, post := getOnePost(w, msg.Data)
			if status != http.StatusOK {
				response.Data = "error"
			} else {
				response.Data = post
			}

		case "getAllPost":
			response.Action = "getAllPost"
			status, posts := getAllPost(w)
			if status != http.StatusOK {
				response.Data = "error"
			} else {
				response.Data = posts
			}

		case "deletePost":
			response.Action = "deletePost"
			if status := deletePost(w, msg.Data); status != http.StatusOK {
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

		fmt.Println("response: ", response)

		responseBytes, _ := json.Marshal(response)
		if err := ws.WriteMessage(websocket.TextMessage, responseBytes); err != nil {
			log.Println(err)
		}
	}

}

func HandleMessages() {
	for {
		msg := <-broadcast
		var message md.MessageChat
		if err := json.Unmarshal(msg, &message); err != nil {
			fmt.Println(err)
			continue
		}
		if recipient, ok := clients[message.ReceiverId]; ok {
			recipient.Send <- msg
		}
	}
}

func writeToClients(client *Client) {
	var response Message
	response.Action = "newMessage"
	for message := range client.Send {
		response.Data = string(message)
		responseBytes, _ := json.Marshal(response)
		err := client.Conn.WriteMessage(websocket.TextMessage, responseBytes)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
