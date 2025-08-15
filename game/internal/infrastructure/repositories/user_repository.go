package repositories

import (
	"context"
	"game/internal/domain/entities"
	"game/internal/domain/repositories"

	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) repositories.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *entities.User) error {
	query := `INSERT INTO users (uuid, login, password_hash, score, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.Exec(ctx, query, user.UUID, user.Login, user.PasswordHash, user.Score, user.CreatedAt, user.UpdatedAt)
	return err
}

func (r *userRepository) Update(ctx context.Context, user *entities.User) error {
	query := `UPDATE users SET login = $1, password_hash = $2, score = $3, updated_at = $4 WHERE uuid = $5`
	_, err := r.db.Exec(ctx, query, user.Login, user.PasswordHash, user.Score, user.UpdatedAt, user.UUID)
	return err
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*entities.User, error) {
	query := `SELECT uuid, login, password_hash, score, created_at, updated_at FROM users WHERE uuid = $1`
	user := &entities.User{}
	err := r.db.QueryRow(ctx, query, id).Scan(&user.UUID, &user.Login, &user.PasswordHash, &user.Score, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetByLogin(ctx context.Context, login string) (*entities.User, error) {
	query := `SELECT uuid, login, password_hash, score, created_at, updated_at FROM users WHERE login = $1`
	user := &entities.User{}
	err := r.db.QueryRow(ctx, query, login).Scan(&user.UUID, &user.Login, &user.PasswordHash, &user.Score, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) IncreaseScore(ctx context.Context, id string) error {
	query := `UPDATE users SET score = score + 1, updated_at = CURRENT_TIMESTAMP WHERE uuid = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
