package handlers

import "errors"

var (
	errInvalidRequest       = errors.New("invalid request body")
	errIncalidLoginPassword = errors.New("invalid login or password")
)
