package model

import "time"

type Comment struct {
	Id         int       `json:"commentID"`
	PostId     int       `json:"postID"`
	UserId     int       `json:"userID"`
	Nickname   string    `json:"nickname"`
	LikedBy    []string  `json:"likedBy"`
	DisLikedBy []string  `json:"dislikedBy"`
	Content    string    `json:"content"`
	NbrLike    int       `json:"nbrLike"`
	NbrDislike int       `json:"nbrDislike"`
	CreateAt   time.Time `json:"createAt"`
}
