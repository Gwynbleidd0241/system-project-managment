package errors

import "errors"

var (
	InvalidInputData  = errors.New("invalid input data")
	TaskAlreadyExists = errors.New("task already exists")
	TaskNotFound      = errors.New("task not found")
)
