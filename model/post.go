package model

import "time"

type Post struct {
	Id         int       `json:"postID"`
	UserId     int       `json:"userID"`
	Categorie  []string  `json:"categorie"`
	Content    string    `json:"content"`
	Img        string    `json:"img"`
	NbrLike    int       `json:"nbrLike"`
	NbrDislike int       `json:"nbrDislike"`
	CreateAt   time.Time `json:"createAt"`
}
