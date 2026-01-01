package errors

import "errors"

var (
	ErrDatabase = errors.New("database operation failed")
	ErrConflict = errors.New("duplicate key or unique constraint violation")
	ErrNotFound = errors.New("record not found")
	ErrParse    = errors.New("failed to parse data")
)
