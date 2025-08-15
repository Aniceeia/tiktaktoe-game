package repositories

import (
	"context"
	"game/internal/domain/entities"
)

type GameRepository interface {
	Create(ctx context.Context, game *entities.Game) error
	Update(ctx context.Context, game *entities.Game) error
	GetByID(ctx context.Context, id string) (*entities.Game, error)
	GetAvailableGames(ctx context.Context) ([]*entities.Game, error)
	GetByPlayerID(ctx context.Context, playerID string) ([]*entities.Game, error)
	GetCompletedByUserID(ctx context.Context, userID string) ([]*entities.Game, error)
	GetLeaderboard(ctx context.Context, limit int) ([]*entities.LeaderboardPlayer, error)
}
