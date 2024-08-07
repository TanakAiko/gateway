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

	status, user := GetUserData(w, sessionID.Value)
	if status != http.StatusOK {
		fmt.Printf("Failed to get user's data: %v\n", err)
		return
	}

	client := &Client{User: user, Conn: ws, Send: make(chan []byte)}
	clients[client.User.Id] = client

	go writeToClients(client)

	status, users := getAllUser(w)
	if status != http.StatusOK {
		fmt.Printf("Failed to get all users: %v\n", err)
		return
	}
	sendAllUser(ws, users, user.Id)

	isAllRead(w, user.Id, ws)

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

		fmt.Println("msg.Action: ", msg.Action)

		switch msg.Action {

		case "logout":
			response.Action = "logout"
			if status := logoutRequest(w, msg.Data); status != http.StatusOK {
				response.Data = "error"
			} else {
				response.Data = "OK"
			}

			delete(clients, user.Id)
			status, users := getAllUser(w)
			if status != http.StatusOK {
				fmt.Printf("Failed to get all users: %v\n", err)
				return
			}
			sendNewUser(users, user.Id)
		case "messageCreate":
			response.Action = "messageCreate"
			if status, jsonMessage := createMessage(w, msg.Data); status != http.StatusOK {
				response.Data = "error"
			} else {
				response.Data = jsonMessage
				broadcast <- []byte(jsonMessage)
			}

		case "messageGets":
			response.Action = "messageGets"
			if status, chats := GetMessages(w); status != http.StatusOK {
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
			if status := createPost(w, msg.Data, client.User.Id); status != http.StatusCreated {
				response.Data = "error"
			} else {
				sendNewPostToAll(w)
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

		case "updatePostLike":
			response.Action = "updatePostLike"
			if status := updatePostLike(w, msg.Data); status != http.StatusOK {
				response.Data = "error"
			} else {
				sendNewLike(w, client.User.Id, msg.Data, "sendNewPostLikes")
				response.Data = "OK"
			}

		case "commentCreate":
			// fmt.Printf("\n\nmsg.Data: %v\n\n", msg.Data)
			response.Action = "commentCreate"
			status, postId := createComment(w, msg.Data)
			if status != http.StatusCreated {
				response.Data = "error"
			} else {
				sendNewCommentToAll(w)
				response.Data = fmt.Sprintf("%d", postId)
			}

		case "updateCommentLike":
			response.Action = "updateCommentLike"
			if status := updateCommentLike(w, msg.Data); status != http.StatusOK {
				response.Data = "error"
			} else {
				sendNewLike(w, client.User.Id, msg.Data, "sendNewCommentLikes")
				response.Data = "OK"
			}

		case "sendNewUsertoAll":
			status, users := getAllUser(w)
			if status != http.StatusOK {
				fmt.Printf("Failed to get all users: %v\n", err)
				return
			}
			sendNewUser(users, user.Id)
			continue

		case "startTyping":
			inTyping("startTyping", msg.Data)
			continue

		case "stopTyping":
			inTyping("stopTyping", msg.Data)
			continue

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

func inTyping(action string, data string) {
	reply := Message{
		Action: action,
		Data:   data,
	}

	var mess md.MessageChat
	if err := json.Unmarshal([]byte(data), &mess); err != nil {
		fmt.Println("Internal server error: " + err.Error())
		return
	}

	if _, ok := clients[mess.ReceiverId]; !ok {
		fmt.Printf("User is not connected")
		return
	}

	responseBytes, _ := json.Marshal(reply)

	err := clients[mess.ReceiverId].Conn.WriteMessage(websocket.TextMessage, responseBytes)
	if err != nil {
		fmt.Println(err)
		return
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

func sendNewPostToAll(w http.ResponseWriter) {
	var replyUpdate Message
	replyUpdate.Action = "updatePosts"
	status, lastPost := getLastPost(w)
	if status == 0 {
		replyUpdate.Data = "error"
	} else {
		replyUpdate.Data = lastPost
	}

	replyUpdate.Data = string(lastPost)

	responseBytes, _ := json.Marshal(replyUpdate)

	for _, client := range clients {
		if err := client.Conn.WriteMessage(websocket.TextMessage, responseBytes); err != nil {
			fmt.Println(err)
			return
		}
	}
}

func sendNewCommentToAll(w http.ResponseWriter) {
	var replyUpdate Message
	replyUpdate.Action = "updateLastComment"
	status, lastComment := getLastComment(w)
	if status == 0 {
		replyUpdate.Data = "error"
	} else {
		replyUpdate.Data = lastComment
	}

	replyUpdate.Data = string(lastComment)

	responseBytes, _ := json.Marshal(replyUpdate)

	for _, client := range clients {
		if err := client.Conn.WriteMessage(websocket.TextMessage, responseBytes); err != nil {
			fmt.Println(err)
			return
		}
	}
}

func sendNewLike(w http.ResponseWriter, userId int, data string, action string) {
	var replyUpdate Message
	replyUpdate.Action = action
	replyUpdate.Data = data
	responseBytes, _ := json.Marshal(replyUpdate)
	for _, client := range clients {
		if client.User.Id == userId {
			continue
		}
		if err := client.Conn.WriteMessage(websocket.TextMessage, responseBytes); err != nil {
			fmt.Println(err)
			return
		}
	}
}

func sendAllUser(ws *websocket.Conn, users []md.User, id int) {
	var newTab []md.User
	var replyUpdate Message
	replyUpdate.Action = "getAllUser"
	for i := 0; i < len(users); i++ {
		if _, ok := clients[users[i].Id]; ok {
			users[i].Status = "online"
		} else {
			users[i].Status = "offline"
		}
		newTab = append(newTab, users[i])
	}

	tabBytes, _ := json.Marshal(newTab)
	replyUpdate.Data = string(tabBytes)

	responseBytes, _ := json.Marshal(replyUpdate)

	if err := ws.WriteMessage(websocket.TextMessage, responseBytes); err != nil {
		fmt.Println(err)
		return
	}
}

func sendNewUser(users []md.User, id int) {
	var newTab []md.User
	var replyUpdate Message
	replyUpdate.Action = "sendNewUsertoAll"

	for i := 0; i < len(users); i++ {
		if _, ok := clients[users[i].Id]; ok {
			users[i].Status = "online"
		} else {
			users[i].Status = "offline"
		}
		newTab = append(newTab, users[i])
	}

	tabBytes, _ := json.Marshal(newTab)
	replyUpdate.Data = string(tabBytes)

	responseBytes, _ := json.Marshal(replyUpdate)

	for _, client := range clients {
		if client.User.Id == id {
			continue
		}
		if err := client.Conn.WriteMessage(websocket.TextMessage, responseBytes); err != nil {
			fmt.Println(err)
			return
		}
	}
}
