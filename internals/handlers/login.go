package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	md "gateway/model"
	"io"
	"net/http"
	"time"
)

type Credential struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

type Session struct {
	Id         string    `json:"sessionID"`
	UserID     int       `json:"userId"`
	Expiration time.Time `json:"sessionExpireTime"`
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
		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("body: ", string(body))

	var credential Credential

	if err = json.Unmarshal(body, &credential); err != nil {
		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	bodyData := md.RequestBody{
		Action: action,
		Body:   credential,
	}

	fmt.Println(bodyData.Body)

	jsonData, err := json.Marshal(bodyData)
	if err != nil {
		http.Error(w, "Internal server error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	req, err := http.NewRequest("POST", URLauth, bytes.NewBuffer(jsonData))
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

	for _, cookie := range resp.Cookies() {
		http.SetCookie(w, cookie)
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Internal server error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	/* var session Session

	if err = json.Unmarshal(responseBody, &session); err != nil {
		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "sessionID",
		Value:    session.Id,
		Expires:  session.Expiration,
		HttpOnly: true,
	})

	fmt.Println("session: ", session) */

	// Transférer les cookies du deuxième serveur au client

	w.WriteHeader(resp.StatusCode)
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBody)
}
