package repositories

import (
	"context"
	"database/sql"
	"game/internal/domain/entities"
	"game/internal/domain/repositories"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) repositories.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *entities.User) error {
	query := `INSERT INTO users (uuid, username, password_hash, score, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.ExecContext(ctx, query, user.UUID, user.Username, user.PasswordHash, user.Score, user.CreatedAt, user.UpdatedAt)
	return err
}

func (r *userRepository) Update(ctx context.Context, user *entities.User) error {
	query := `UPDATE users SET username = $1, password_hash = $2, score = $3, updated_at = $4 WHERE uuid = $5`
	_, err := r.db.ExecContext(ctx, query, user.Username, user.PasswordHash, user.Score, user.UpdatedAt, user.UUID)
	return err
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*entities.User, error) {
	query := `SELECT uuid, username, password_hash, score, created_at, updated_at FROM users WHERE uuid = $1`
	user := &entities.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&user.UUID, &user.Username, &user.PasswordHash, &user.Score, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*entities.User, error) {
	query := `SELECT uuid, username, password_hash, score, created_at, updated_at FROM users WHERE username = $1`
	user := &entities.User{}
	err := r.db.QueryRowContext(ctx, query, username).Scan(&user.UUID, &user.Username, &user.PasswordHash, &user.Score, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) IncreaseScore(ctx context.Context, id string) error {
	query := `UPDATE users SET score = score + 10, updated_at = CURRENT_TIMESTAMP WHERE uuid = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
