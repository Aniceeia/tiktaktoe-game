package services

import "errors"

var (
	errPVPJoinOnly    = errors.New("can only join PvP games")
	errGameFull       = errors.New("game is full")
	errUserExist      = errors.New("user already exists")
	errBcryptGenerate = errors.New("hash error")
	errCreateUser     = errors.New("user creation error")
	errLogin          = errors.New("invalid login")
	errPassword       = errors.New("invalid password")
)
