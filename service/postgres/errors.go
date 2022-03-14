package postgres

import "errors"

var (
	ErrAlreadyExists    = errors.New("already exists")
	ErrResourceNotFound = errors.New("resource not found")
)
