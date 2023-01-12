package model

type User struct {
	ID		 string
	Username string `json:"login" validate:"required,ascii"`
	Password string `json:"password" validate:"required,ascii"`
}




