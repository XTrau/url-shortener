package apperrors

import "errors"

var (
	ErrUrlNotFound   = errors.New("Url not found")
	ErrAlreadyExists = errors.New("Already exists")
)
