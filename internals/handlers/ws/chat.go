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

func createMessage(w http.ResponseWriter, data string) int {
	var mess md.MessageChat
	if err := json.Unmarshal([]byte(data), &mess); err != nil {
		fmt.Println("Internal server error: " + err.Error())
		return 0
	}

	bodyData := md.RequestBody{
		Action: "createChat",
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
		http.Error(w, "Internal server error: "+err.Error(), http.StatusInternalServerError)
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
