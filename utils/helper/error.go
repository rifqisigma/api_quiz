package helper

import "errors"

var (
	ErrUserNotFound = errors.New("uset not found")
	ErrServerError  = errors.New("server error")
	ErrInvalidEmail = errors.New("invalid email")
)
