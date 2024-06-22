package handlers

import (
	"encoding/json"
	"fmt"
	"gateway/internals/tools"
	md "gateway/model"
	"io"
	"net/http"
)

type Credential struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	const action = "login"
	fmt.Println("action: ", action)

	if r.Method != http.MethodPost {
		http.Error(w, "methode not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Bad request: " + err.Error())
		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	var credential Credential

	if err = json.Unmarshal(body, &credential); err != nil {
		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	bodyData := md.RequestBody{
		Action: action,
		Body:   credential,
	}

	resp, err := tools.SendRequest(w, bodyData, "POST", URLauth)
	if err != nil {
		fmt.Println("Internal server error: " + err.Error())
		http.Error(w, "Internal server error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	for _, cookie := range resp.Cookies() {
		fmt.Println("cookie : ", cookie)
		http.SetCookie(w, cookie)
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Internal server error: " + err.Error())
		http.Error(w, "Internal server error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Transférer les cookies du deuxième serveur au client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(responseBody)
}
