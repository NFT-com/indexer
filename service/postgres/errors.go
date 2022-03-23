package postgres

import "errors"

var (
	errAlreadyExists    = errors.New("already exists")
	errResourceNotFound = errors.New("resource not found")
)
