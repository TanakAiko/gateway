package handlers

import (
	"encoding/json"
	"fmt"
	conf "gateway/config"
	"gateway/internals/tools"
	md "gateway/model"
	"io"
	"net/http"
)

func GetUserData(w http.ResponseWriter, sessionID string) (int, md.User) {
	var user md.User
	user.SessionID = sessionID
	bodyData := md.RequestBody{
		Action: "getUserData",
		Body:   user,
	}

	fmt.Println("action: ", bodyData.Action)

	resp, err := tools.SendRequest(w, bodyData, "POST", conf.URLauth)
	if err != nil {
		fmt.Println("Internal server error: " + err.Error())
		return 0, user
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Internal server error: " + err.Error())
		return 0, user
	}

	if err = json.Unmarshal(responseBody, &user); err != nil {
		return 0, user
	}

	return resp.StatusCode, user
}
