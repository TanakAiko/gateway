package model

import "time"

type MessageChat struct {
	Id             int       `json:"messageID"`
	SenderId       int       `json:"senderID"`
	ReceiverId     int       `json:"receiverID"`
	Content        string    `json:"content"`
	StatusReceived bool      `json:"statusReceived"`
	StatusRead     bool      `json:"statusRead"`
	CreateAt       time.Time `json:"createAT"`
}
