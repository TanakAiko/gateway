package handlers

import (
	"encoding/json"
	"fmt"
	conf "gateway/config"
	"gateway/internals/tools"
	md "gateway/model"
	"net/http"
)

/* func getAllComment(w http.ResponseWriter, postIdStr string) (int, string) {
	var comment md.Comment

	postId, err := strconv.Atoi(postIdStr)
	if err != nil {
		fmt.Println("Internal server error: " + err.Error())
		return 0, ""
	}
	comment.PostId = postId

	bodyData := md.RequestBody{
		Action: "getAll",
		Body:   comment,
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
} */

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
