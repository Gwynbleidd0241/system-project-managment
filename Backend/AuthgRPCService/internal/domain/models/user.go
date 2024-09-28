package models

type User struct {
	ID       string `swaggerignore:"true"`
	Email    string `validate:"required,email" json:"email"`
	PassHash string `validate:"required,min=5,max=70" json:"password"`
}
