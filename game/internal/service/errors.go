package service

import "errors"

var (
	errGameNotFound  = errors.New("game not found")
	errInvalidBoard  = errors.New("invalid board structure")
	errInvalidMove   = errors.New("invalid move detected")
	errGameCompleted = errors.New("game is already completed")
	errWrongTurn     = errors.New("not your turn")
	errGameFull      = errors.New("game is full")
	errInvalidMode   = errors.New("invalid game mode")
	//errNotImplInServ = errors.New("not implemented in service layer")
	ErrUserExists = errors.New("user already exists")
)
