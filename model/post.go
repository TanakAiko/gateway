package model

import "time"

type Post struct {
	Id         int       `json:"postID"`
	UserId     int       `json:"userID"`
	Nickname   string    `json:"nickname"`
	Categorie  []string  `json:"categorie"`
	LikedBy    []string  `json:"likedBy"`
	DislikedBy []string  `json:"dislikedBy"`
	Content    string    `json:"content"`
	Img        string    `json:"img"`
	ImgBase64  string    `json:"imgBase64"`
	NbrLike    int       `json:"nbrLike"`
	NbrDislike int       `json:"nbrDislike"`
	Comments   string    `json:"comments"`
	CreateAt   time.Time `json:"createAt"`
}
