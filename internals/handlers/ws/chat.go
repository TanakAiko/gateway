package handlers

import (
	"encoding/json"
	"fmt"
	conf "gateway/config"
	"gateway/internals/tools"
	md "gateway/model"
	"io"
	"net/http"

	"github.com/gorilla/websocket"
)

func createMessage(w http.ResponseWriter, data string) (int, string) {
	var mess md.MessageChat
	if err := json.Unmarshal([]byte(data), &mess); err != nil {
		fmt.Println("Internal server error: " + err.Error())
		return 0, ""
	}

	bodyData := md.RequestBody{
		Action: "createChat",
		Body:   mess,
	}

	resp, err := tools.SendRequest(w, bodyData, "POST", conf.URLChat)
	if err != nil {
		fmt.Println("Internal server error: " + err.Error())
		return 0, ""
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Internal server error: " + err.Error())
		return 0, ""
	}

	return resp.StatusCode, string(responseBody)
}

func getMessages(w http.ResponseWriter) (int, string) {
	bodyData := md.RequestBody{
		Action: "getChats",
	}

	resp, err := tools.SendRequest(w, bodyData, "POST", conf.URLChat)
	if err != nil {
		fmt.Println("Internal server error: " + err.Error())
		return 0, ""
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Internal server error: " + err.Error())
		return 0, ""
	}

	return resp.StatusCode, string(responseBody)
}

func setStatusReceived(w http.ResponseWriter, data string) int {
	var mess md.MessageChat
	if err := json.Unmarshal([]byte(data), &mess); err != nil {
		fmt.Println("Internal server error: " + err.Error())
		return 0
	}

	bodyData := md.RequestBody{
		Action: "updateStatusReceived",
		Body:   mess,
	}

	resp, err := tools.SendRequest(w, bodyData, "POST", conf.URLChat)
	if err != nil {
		fmt.Println("Internal server error: " + err.Error())
		return 0
	}
	defer resp.Body.Close()

	return resp.StatusCode
}

func setStatusRead(w http.ResponseWriter, data string) int {
	var mess md.MessageChat
	if err := json.Unmarshal([]byte(data), &mess); err != nil {
		fmt.Println("Internal server error: " + err.Error())
		return 0
	}

	bodyData := md.RequestBody{
		Action: "updateStatusRead",
		Body:   mess,
	}

	resp, err := tools.SendRequest(w, bodyData, "POST", conf.URLChat)
	if err != nil {
		fmt.Println("Internal server error: " + err.Error())
		return 0
	}
	defer resp.Body.Close()

	return resp.StatusCode
}

func isAllRead(w http.ResponseWriter, userId int, ws *websocket.Conn) {
	var replyUpdate Message
	replyUpdate.Action = "isAllRead"

	status, messagesString := getMessages(w)
	if status != http.StatusOK {
		fmt.Println("Internal server error (isAllRead): try to get message")
		return
	}
	var messages []md.MessageChat
	if err := json.Unmarshal([]byte(messagesString), &messages); err != nil {
		fmt.Println("Internal server error (isAllRead): try to Unmarshal")
		return
	}

	var idSender []int
	for _, msg := range messages {
		if msg.ReceiverId == userId && !msg.StatusRead {
			idSender = append(idSender, msg.SenderId)
		}
	}

	idSenderString, err := json.Marshal(idSender)
	if err != nil {
		fmt.Println("Internal server error (isAllRead): try to Unmarshal")
		return
	}

	replyUpdate.Data = string(idSenderString)

	responseBytes, _ := json.Marshal(replyUpdate)

	if err := ws.WriteMessage(websocket.TextMessage, responseBytes); err != nil {
		fmt.Println(err)
		return
	}
}
