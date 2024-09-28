package models

type User struct {
	ID       string `swaggerignore:"true"`
	Email    string `json:"email"`
	PassHash string `json:"password"`
}
