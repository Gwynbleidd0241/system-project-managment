package errors

import "errors"

var (
	InvalidInputData  = errors.New("invalid input data")
	UserAlreadyExists = errors.New("user already exists")
	UserNotFound      = errors.New("user not found")
)
