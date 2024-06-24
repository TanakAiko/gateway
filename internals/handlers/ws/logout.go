package handlers

import (
	"fmt"
	conf "gateway/config"
	"gateway/internals/tools"
	md "gateway/model"
	"net/http"
)

func logoutRequest(w http.ResponseWriter, sessionID string) int {
	var user md.User
	user.SessionID = sessionID
	bodyData := md.RequestBody{
		Action: "logout",
		Body:   user,
	}

	resp, err := tools.SendRequest(w, bodyData, "POST", conf.URLauth)
	if err != nil {
		fmt.Println("Internal server error: " + err.Error())
		return 0
	}
	defer resp.Body.Close()

	return resp.StatusCode
}
