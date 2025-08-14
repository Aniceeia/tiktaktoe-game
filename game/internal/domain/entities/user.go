package entities

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UUID         string    `json:"uuid"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"password_hash"`
	Score        int       `json:"score"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func NewUser(username, passwordHash string) *User {
	now := time.Now()
	return &User{
		UUID:         generateUserID(),
		Username:     username,
		PasswordHash: passwordHash,
		Score:        0,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

func generateUserID() string {
	return uuid.New().String()
}

// IncreaseScore увеличивает счет пользователя
func (u *User) IncreaseScore(points int) {
	u.Score += points
	u.UpdatedAt = time.Now()
}

// UpdateScore обновляет счет пользователя
func (u *User) UpdateScore(score int) {
	u.Score = score
	u.UpdatedAt = time.Now()
}
