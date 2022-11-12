package database

import "errors"

var (
	ErrTokenNotFound = errors.New("no matching rows for the given token")
	ErrUrlNotFound   = errors.New("no matching rows for the given url")
)
