package entities

import "errors"

var (
	ErrGameNotFound       = errors.New("game not found")
	ErrInvalidBoard       = errors.New("invalid board structure")
	ErrInvalidMove        = errors.New("invalid move detected")
	ErrWrongTurn          = errors.New("not your turn")
	ErrGameFull           = errors.New("game is full")
	ErrInvalidMode        = errors.New("invalid game mode")
	ErrUserExists         = errors.New("user already exists")
	ErrGameNotRunning     = errors.New("game is not in progress")
	ErrInvalidPosition    = errors.New("position already occupied")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrForbidden          = errors.New("forbidden")
)
