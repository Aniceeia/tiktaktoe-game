package entities

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UUID         string    `json:"uuid" db:"uuid"`
	Login        string    `json:"login" db:"login"`
	PasswordHash string    `json:"password_hash" db:"password_hash"`
	Score        int       `json:"score" db:"score"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

func NewUser(login, passwordHash string) *User {
	now := time.Now()
	return &User{
		UUID:         generateUserID(),
		Login:        login,
		PasswordHash: passwordHash,
		Score:        0,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

func generateUserID() string {
	return uuid.New().String()
}

func (u *User) IncreaseScore(points int) {
	u.Score += points
	u.UpdatedAt = time.Now()
}

func (u *User) UpdateScore(score int) {
	u.Score = score
	u.UpdatedAt = time.Now()
}
