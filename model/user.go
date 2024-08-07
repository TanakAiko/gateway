package model

type User struct {
	Id        int    `json:"userId"`
	Nickname  string `json:"nickname"`
	Age       int    `json:"age"`
	Gender    string `json:"gender"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	SessionID string `json:"sessionID"`
	Status    string `json:"status"`
}
