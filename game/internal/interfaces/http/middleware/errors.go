package middleware

import "errors"

var (
	errUnauthorized       = errors.New("authorization header is required")
	errInvalidHeader      = errors.New("invalid bearer token format")
	errUserNotFound       = "invalid user"
	errInvalidContentType = errors.New("content-type must be application/json")
	errInvalidJSON        = errors.New("invalid JSON format")
	errValidationFailed   = errors.New("validation failed")
)
