package service

import (
	"database/sql"
	"encoding/json"

	"github.com/google/uuid"
)

const (
	empty = 0
	x     = 1
	o     = 2
)

// Game struct for id, board and status of the game
type Game struct {
	ID       string
	Board    [3][3]int
	Status   string
	Player1  string // user ID or username
	Player2  string // user ID or username, or "AI" for PvE
	NextTurn string // user ID or "AI"
	Mode     string // "pvp" or "pve"
}

func (g *Game) updateGameStatus() {
	for i := range 3 {
		if g.Board[i][0] != empty &&
			g.Board[i][0] == g.Board[i][1] &&
			g.Board[i][1] == g.Board[i][2] && g.Board[i][2] == x {
			g.Status = StatusPlayer1Won + ":" + g.Player1
			return
		}
	}

	for j := range 3 {
		if g.Board[0][j] != empty &&
			g.Board[0][j] == g.Board[1][j] &&
			g.Board[1][j] == g.Board[2][j] && g.Board[2][j] == o {
			g.Status = StatusPlayer2Won + ":" + g.Player2
			return
		}
	}

	if g.Board[0][0] != empty &&
		g.Board[0][0] == g.Board[1][1] &&
		g.Board[1][1] == g.Board[2][2] && g.Board[2][2] == x {
		g.Status = StatusPlayer1Won + ":" + g.Player1
		return
	}
	if g.Board[0][2] != empty &&
		g.Board[0][2] == g.Board[1][1] &&
		g.Board[1][1] == g.Board[2][0] && g.Board[2][0] == o {
		g.Status = StatusPlayer2Won + ":" + g.Player2
		return
	}

	isFull := true
	for i := range 3 {
		for j := range 3 {
			if g.Board[i][j] == empty {
				isFull = false
				break
			}
		}
	}

	if isFull {
		g.Status = StatusDraw
	} else {
		g.Status = StatusWaiting
	}
}

func GenerateID() string {
	return uuid.New().String()
}

// Добавляем методы для сканирования из БД
func (g *Game) Scan(row *sql.Row) error {
	var boardJSON []byte
	err := row.Scan(
		&g.ID,
		&boardJSON,
		&g.Status,
		&g.Player1,
		&g.Player2,
		&g.NextTurn,
		&g.Mode,
	)
	if err != nil {
		return err
	}
	return json.Unmarshal(boardJSON, &g.Board)
}

func (g *Game) ScanRows(rows *sql.Rows) error {
	var boardJSON []byte
	err := rows.Scan(
		&g.ID,
		&boardJSON,
		&g.Status,
		&g.Player1,
		&g.Player2,
		&g.NextTurn,
		&g.Mode,
	)
	if err != nil {
		return err
	}
	return json.Unmarshal(boardJSON, &g.Board)
}
