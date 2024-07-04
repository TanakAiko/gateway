package model

import "time"

type Post struct {
	Id         int       `json:"postID"`
	UserId     int       `json:"userID"`
	Nickname   string    `json:"nickname"`
	Categorie  []string  `json:"categorie"`
	Content    string    `json:"content"`
	Img        string    `json:"img"`
	ImgBase64  string    `json:"imgBase64"`
	NbrLike    int       `json:"nbrLike"`
	NbrDislike int       `json:"nbrDislike"`
	CreateAt   time.Time `json:"createAt"`
}
