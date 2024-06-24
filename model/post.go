package model

import "time"

type Post struct {
	Id        int       `json:"postID"`
	UserId    int       `json:"userID"`
	Categorie []string  `json:"categorie"`
	Content   string    `json:"content"`
	CreateAt  time.Time `json:"createAt"`
}
