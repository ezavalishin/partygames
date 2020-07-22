package models

type ErrorLog struct {
	BaseModel
	User    *User
	UserId  int
	Ua      *string `json:"ua"`
	Payload *string `json:"payload"`
}
