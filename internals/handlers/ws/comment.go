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

func createComment(w http.ResponseWriter, data string) (int, int) {
	var comment md.Comment

	if err := json.Unmarshal([]byte(data), &comment); err != nil {
		fmt.Println("Internal server error: " + err.Error())
		return 0, 0
	}

	bodyData := md.RequestBody{
		Action: "createComment",
		Body:   comment,
	}

	resp, err := tools.SendRequest(w, bodyData, "POST", conf.URLComment)
	if err != nil {
		fmt.Println("Internal server error: " + err.Error())
		return 0, 0
	}
	defer resp.Body.Close()

	return resp.StatusCode, comment.PostId
}

func getLastComment(w http.ResponseWriter) (int, string) {
	bodyData := md.RequestBody{
		Action: "getLastComment",
	}

	resp, err := tools.SendRequest(w, bodyData, "POST", conf.URLComment)
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

func updateCommentLike(w http.ResponseWriter, data string) int {
	var comment md.Comment
	if err := json.Unmarshal([]byte(data), &comment); err != nil {
		fmt.Println("Internal server error: " + err.Error())
		return 0
	}

	bodyData := md.RequestBody{
		Action: "updateLike",
		Body:   comment,
	}

	resp, err := tools.SendRequest(w, bodyData, "POST", conf.URLComment)
	if err != nil {
		fmt.Println("Internal server error: " + err.Error())
		return 0
	}
	defer resp.Body.Close()

	return resp.StatusCode
}
