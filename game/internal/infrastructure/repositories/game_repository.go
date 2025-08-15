package repositories

import (
	"context"
	"encoding/json"
	"game/internal/domain/entities"
	"game/internal/domain/repositories"

	"github.com/jackc/pgx/v5/pgxpool"
)

type gameRepository struct {
	db *pgxpool.Pool
}

func NewGameRepository(db *pgxpool.Pool) repositories.GameRepository {
	return &gameRepository{db: db}
}

func (r *gameRepository) Create(ctx context.Context, game *entities.Game) error {
	boardJSON, err := json.Marshal(game.Board)
	if err != nil {
		return err
	}

	query := `INSERT INTO games (id, board, status, player1_id, player2_id, next_turn_id, mode, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err = r.db.Exec(ctx, query, game.ID, boardJSON, game.Status, game.Player1ID, game.Player2ID, game.NextTurnID, game.Mode, game.CreatedAt, game.UpdatedAt)
	return err
}

func (r *gameRepository) Update(ctx context.Context, game *entities.Game) error {
	boardJSON, err := json.Marshal(game.Board)
	if err != nil {
		return err
	}

	query := `UPDATE games SET board = $1, status = $2, player1_id = $3, player2_id = $4, next_turn_id = $5, mode = $6, updated_at = $7 WHERE id = $8`
	_, err = r.db.Exec(ctx, query, boardJSON, game.Status, game.Player1ID, game.Player2ID, game.NextTurnID, game.Mode, game.UpdatedAt, game.ID)
	return err
}

func (r *gameRepository) GetByID(ctx context.Context, id string) (*entities.Game, error) {
	query := `SELECT id, board, status, player1_id, player2_id, next_turn_id, mode, created_at, updated_at FROM games WHERE id = $1`

	var boardJSON []byte
	var player2 *string

	game := &entities.Game{}
	err := r.db.QueryRow(ctx, query, id).Scan(&game.ID, &boardJSON, &game.Status, &game.Player1ID, &player2, &game.NextTurnID, &game.Mode, &game.CreatedAt, &game.UpdatedAt)
	if err != nil {
		return nil, err
	}

	if player2 != nil {
		game.Player2ID = *player2
	}

	err = json.Unmarshal(boardJSON, &game.Board)
	if err != nil {
		return nil, err
	}

	return game, nil
}

func (r *gameRepository) GetAvailableGames(ctx context.Context) ([]*entities.Game, error) {
	query := `SELECT id, board, status, player1_id, player2_id, next_turn_id, mode, created_at, updated_at FROM games WHERE status = 'waiting' AND mode = 'pvp' ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var games []*entities.Game
	for rows.Next() {
		var boardJSON []byte
		var player2 *string

		game := &entities.Game{}
		err := rows.Scan(&game.ID, &boardJSON, &game.Status, &game.Player1ID, &player2, &game.NextTurnID, &game.Mode, &game.CreatedAt, &game.UpdatedAt)
		if err != nil {
			return nil, err
		}

		if player2 != nil {
			game.Player2ID = *player2
		}

		err = json.Unmarshal(boardJSON, &game.Board)
		if err != nil {
			return nil, err
		}

		games = append(games, game)
	}

	return games, nil
}

func (r *gameRepository) GetByPlayerID(ctx context.Context, playerID string) ([]*entities.Game, error) {
	query := `SELECT id, board, status, player1_id, player2_id, next_turn_id, mode, created_at, updated_at FROM games WHERE player1_id = $1 OR player2_id = $1 ORDER BY updated_at DESC`

	rows, err := r.db.Query(ctx, query, playerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var games []*entities.Game
	for rows.Next() {
		var boardJSON []byte
		var player2 *string

		game := &entities.Game{}
		err := rows.Scan(&game.ID, &boardJSON, &game.Status, &game.Player1ID, &player2, &game.NextTurnID, &game.Mode, &game.CreatedAt, &game.UpdatedAt)
		if err != nil {
			return nil, err
		}

		if player2 != nil {
			game.Player2ID = *player2
		}

		err = json.Unmarshal(boardJSON, &game.Board)
		if err != nil {
			return nil, err
		}

		games = append(games, game)
	}

	return games, nil
}
