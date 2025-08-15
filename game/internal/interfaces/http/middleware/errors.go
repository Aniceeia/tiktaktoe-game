package middleware

import "errors"

var (
	errInvalidHeader = errors.New("invalid bearer token format")
)
