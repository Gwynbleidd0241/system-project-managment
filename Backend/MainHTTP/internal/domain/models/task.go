package models

type Task struct {
	ID          string `swagignore:"true"`
	UserEmail   string `json:"email"`
	Name        string `json:"title"`
	Description string `json:"body"`
	IsDone      bool   `json:"is_done"`
	Token       string `json:"token"`
}
