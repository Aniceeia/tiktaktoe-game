package service

import (
	"database/sql"
	"encoding/json"
	"errors"
)

type GameRepository interface {
	Save(game *Game) error
	Get(id string) (*Game, error)
	Exists(id string) bool
	GetAllGames() ([]*Game, error)
}

type PostgresGameStorage struct {
	db *sql.DB
}

func NewPostgresGameStorage(db *sql.DB) GameRepository {
	return &PostgresGameStorage{db: db}
}

func (s *PostgresGameStorage) Save(game *Game) error {
	boardJSON, err := json.Marshal(game.Board)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(`
		INSERT INTO games (id, board, status, player1, player2, next_turn, mode)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (id) DO UPDATE 
		SET board = EXCLUDED.board, 
		    status = EXCLUDED.status,
		    player1 = EXCLUDED.player1,
		    player2 = EXCLUDED.player2,
		    next_turn = EXCLUDED.next_turn,
		    mode = EXCLUDED.mode
	`,
		game.ID,
		boardJSON,
		game.Status,
		game.Player1,
		game.Player2,
		game.NextTurn,
		game.Mode,
	)
	return err
}

func (s *PostgresGameStorage) Get(id string) (*Game, error) {
	game := &Game{}
	err := game.Scan(s.db.QueryRow(`
		SELECT id, board, status, player1, player2, next_turn, mode 
		FROM games 
		WHERE id = $1
	`, id))

	if errors.Is(err, sql.ErrNoRows) {
		return nil, errGameNotFound
	}
	return game, err
}

func (s *PostgresGameStorage) Exists(id string) bool {
	var exists bool
	err := s.db.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM games WHERE id = $1)
	`, id).Scan(&exists)
	return err == nil && exists
}

func (s *PostgresGameStorage) GetAllGames() ([]*Game, error) {
	rows, err := s.db.Query(`
		SELECT id, board, status, player1, player2, next_turn, mode 
		FROM games
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var games []*Game
	for rows.Next() {
		game := &Game{}
		if err := game.ScanRows(rows); err != nil {
			return nil, err
		}
		games = append(games, game)
	}
	return games, nil
}
