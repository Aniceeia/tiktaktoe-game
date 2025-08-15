package services

import "errors"

var (
	errBcryptGenerate = errors.New("hash error")
	errCreateUser     = errors.New("user creation error")
)
