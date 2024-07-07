package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	conf "gateway/config"
	md "gateway/model"
	"io"
	"net/http"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	const action = "register"

	fmt.Println("action: ", action)

	if r.Method != http.MethodPost {
		http.Error(w, "methode not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	var user md.User

	if err = json.Unmarshal(body, &user); err != nil {
		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	bodyData := md.RequestBody{
		Action: action,
		Body:   user,
	}

	jsonData, err := json.Marshal(bodyData)
	if err != nil {
		http.Error(w, "Internal server error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	req, err := http.NewRequest("POST", conf.URLauth, bytes.NewBuffer(jsonData))
	if err != nil {
		http.Error(w, "Internal server error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Internal server error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Internal server error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(resp.StatusCode)
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBody)
}
