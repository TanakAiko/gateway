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

func createPost(w http.ResponseWriter, data string) int {
	var post md.Post
	if err := json.Unmarshal([]byte(data), &post); err != nil {
		fmt.Println("Internal server error: " + err.Error())
		return 0
	}

	bodyData := md.RequestBody{
		Action: "createPost",
		Body:   post,
	}

	resp, err := tools.SendRequest(w, bodyData, "POST", conf.URLPost)
	if err != nil {
		fmt.Println("Internal server error: " + err.Error())
		return 0
	}
	defer resp.Body.Close()

	return resp.StatusCode
}

func getOnePost(w http.ResponseWriter, data string) (int, string) {
	var post md.Post
	if err := json.Unmarshal([]byte(data), &post); err != nil {
		fmt.Println("Internal server error: " + err.Error())
		return 0, ""
	}

	bodyData := md.RequestBody{
		Action: "getOne",
		Body:   post,
	}

	resp, err := tools.SendRequest(w, bodyData, "POST", conf.URLPost)
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

func getAllPost(w http.ResponseWriter) (int, string) {
	bodyData := md.RequestBody{
		Action: "getAll",
	}

	resp, err := tools.SendRequest(w, bodyData, "POST", conf.URLPost)
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

func deletePost(w http.ResponseWriter, data string) int {
	var post md.Post
	if err := json.Unmarshal([]byte(data), &post); err != nil {
		fmt.Println("Internal server error: " + err.Error())
		return 0
	}

	bodyData := md.RequestBody{
		Action: "delete",
		Body:   post,
	}

	resp, err := tools.SendRequest(w, bodyData, "POST", conf.URLPost)
	if err != nil {
		fmt.Println("Internal server error: " + err.Error())
		return 0
	}
	defer resp.Body.Close()

	return resp.StatusCode

}
