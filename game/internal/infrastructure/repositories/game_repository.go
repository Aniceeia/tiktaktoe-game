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

func (r *gameRepository) GetCompletedByUserID(ctx context.Context, userID string) ([]*entities.Game, error) {
	query := `
        SELECT id, board, status, player1_id, player2_id, next_turn_id, mode, created_at, updated_at
        FROM games
        WHERE (
            status = 'draw' AND (player1_id = $1 OR player2_id = $1)
        ) OR (
            status = 'player1_won' AND player1_id = $1
        ) OR (
            status = 'player2_won' AND player2_id = $1
        )
        ORDER BY updated_at DESC`

	rows, err := r.db.Query(ctx, query, userID)
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

		if err := json.Unmarshal(boardJSON, &game.Board); err != nil {
			return nil, err
		}

		games = append(games, game)
	}

	return games, nil
}

func (r *gameRepository) GetLeaderboard(ctx context.Context, limit int) ([]*entities.LeaderboardPlayer, error) {
	query := `
        WITH player_stats AS (
            SELECT 
                u.uuid,
                u.login,
                COUNT(CASE WHEN g.status = 'player1_won' AND g.player1_id = u.uuid THEN 1 END) +
                COUNT(CASE WHEN g.status = 'player2_won' AND g.player2_id = u.uuid THEN 1 END) as wins,
                COUNT(CASE WHEN g.status = 'player1_won' AND g.player2_id = u.uuid THEN 1 END) +
                COUNT(CASE WHEN g.status = 'player2_won' AND g.player1_id = u.uuid THEN 1 END) as losses,
                COUNT(CASE WHEN g.status = 'draw' AND (g.player1_id = u.uuid OR g.player2_id = u.uuid) THEN 1 END) as draws
            FROM users u
            LEFT JOIN games g ON (g.player1_id = u.uuid OR g.player2_id = u.uuid) 
                AND g.status IN ('player1_won', 'player2_won', 'draw')
            GROUP BY u.uuid, u.login
        )
        SELECT 
            uuid,
            login,
            CASE 
                WHEN (wins + losses + draws) = 0 THEN 0.0
                ELSE ROUND((wins::float / (wins + losses + draws)::float)::numeric, 3)
            END as win_ratio,
            wins,
            losses,
            draws,
            (wins + losses + draws) as total_games
        FROM player_stats
        WHERE (wins + losses + draws) > 0
        ORDER BY win_ratio DESC, wins DESC
        LIMIT $1`

	rows, err := r.db.Query(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var players []*entities.LeaderboardPlayer
	for rows.Next() {
		player := &entities.LeaderboardPlayer{}
		err := rows.Scan(&player.UUID, &player.Login, &player.WinRatio, &player.Wins, &player.Losses, &player.Draws, &player.TotalGames)
		if err != nil {
			return nil, err
		}
		players = append(players, player)
	}

	return players, nil
}
