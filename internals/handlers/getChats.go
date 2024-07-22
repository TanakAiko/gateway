package handlers

import (
	"fmt"
	ws "gateway/internals/handlers/ws"
	"net/http"
)

func GetAllChatHandler(w http.ResponseWriter, r *http.Request) {
	status, chats := ws.GetMessages(w)
	if status != http.StatusOK {
		fmt.Printf("Failed to get user's message\n")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(chats))
}
