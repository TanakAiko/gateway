package model

type RequestBody struct {
	Action string `json:"action"`
	Body   any    `json:"body"`
}
