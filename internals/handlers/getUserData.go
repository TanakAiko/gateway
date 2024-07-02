package handlers

import (
	"encoding/json"
	"fmt"
	ws "gateway/internals/handlers/ws"
	"net/http"
)

func GetUserDataHandler(w http.ResponseWriter, r *http.Request) {
	sessionID, err := r.Cookie("sessionID")
	if err != nil {
		fmt.Printf("Failed to get cookie: %v\n", err)
		return
	}

	status, user := ws.GetUserData(w, sessionID.Value)
	if status != http.StatusOK {
		fmt.Printf("Failed to get user's data: %v\n", err)
		return
	}

	jsonData, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "Internal server error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(jsonData)
}
