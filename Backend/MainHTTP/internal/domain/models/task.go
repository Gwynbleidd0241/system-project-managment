package models

type Task struct {
	ID          string `swagignore:"true"`
	UserEmail   string `json:"email"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
