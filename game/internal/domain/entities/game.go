package entities

import (
	"time"

	"github.com/google/uuid"
)

type GameStatus string
type GameMode string

const (
	StatusWaiting     GameStatus = "waiting"
	StatusPlayer1Turn GameStatus = "player1_turn"
	StatusPlayer2Turn GameStatus = "player2_turn"
	StatusDraw        GameStatus = "draw"
	StatusPlayer1Won  GameStatus = "player1_won"
	StatusPlayer2Won  GameStatus = "player2_won"
	StatusAITurn      GameStatus = "ai_turn"

	ModePvP GameMode = "pvp"
	ModePvE GameMode = "pve"
)

type Game struct {
	ID         string     `json:"id" db:"id"`
	Board      [3][3]int  `json:"board" db:"board"`
	Status     GameStatus `json:"status" db:"status"`
	Player1ID  string     `json:"player1_id" db:"player1_id"`
	Player2ID  string     `json:"player2_id" db:"player2_id"`
	NextTurnID string     `json:"next_turn_id" db:"next_turn_id"`
	Mode       GameMode   `json:"mode" db:"mode"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at" db:"updated_at"`
}

type LeaderboardPlayer struct {
	UUID       string  `json:"uuid" db:"uuid"`
	Login      string  `json:"login" db:"login"`
	WinRatio   float64 `json:"win_ratio" db:"win_ratio"`
	Wins       int     `json:"wins" db:"wins"`
	Losses     int     `json:"losses" db:"losses"`
	Draws      int     `json:"draws" db:"draws"`
	TotalGames int     `json:"total_games" db:"total_games"`
}

// 1 - player "X" (Player1)
// 2 - player "O" (Player2/AI)

func NewGame(creator string, mode GameMode) *Game {
	now := time.Now()
	status := StatusWaiting
	if mode == ModePvE {
		status = StatusPlayer1Turn
	}

	return &Game{
		ID:         generateID(),
		Board:      [3][3]int{},
		Status:     status,
		Player1ID:  creator,
		Player2ID:  "",
		NextTurnID: creator,
		Mode:       mode,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

func generateID() string {
	return uuid.New().String()
}

func (g *Game) checkErrors(playerID string, row, col int) error {
	if g.Status != StatusPlayer1Turn && g.Status != StatusPlayer2Turn {
		return ErrGameNotRunning
	}
	if g.NextTurnID != playerID {
		return ErrWrongTurn
	}
	if row < 0 || row > 2 || col < 0 || col > 2 {
		return ErrInvalidMove
	}
	if g.Board[row][col] != 0 {
		return ErrInvalidPosition
	}
	return nil
}

func (g *Game) writeSymbol(playerID string, row, col int) {
	symbol := 1 // X for Player1
	if g.Player2ID == playerID {
		symbol = 2 // O for Player2
	}
	g.Board[row][col] = symbol
}

func (g *Game) MakeMove(playerID string, row, col int) error {
	if err := g.checkErrors(playerID, row, col); err != nil {
		return err
	}

	g.writeSymbol(playerID, row, col)
	g.UpdateStatus()
	g.SwitchPlayer()
	g.UpdateNextPlayer()
	g.UpdatedAt = time.Now()

	return nil
}

func (g *Game) UpdateStatus() {
	if g.isWinner(1) {
		g.Status = StatusPlayer1Won
		return
	}
	if g.isWinner(2) {
		g.Status = StatusPlayer2Won
		return
	}
	if g.isBoardFull() {
		g.Status = StatusDraw
		return
	}
}

func (g *Game) isWinner(symbol int) bool {
	for i := 0; i < 3; i++ {
		if g.Board[i][0] == symbol && g.Board[i][1] == symbol && g.Board[i][2] == symbol {
			return true
		}
	}

	for j := 0; j < 3; j++ {
		if g.Board[0][j] == symbol && g.Board[1][j] == symbol && g.Board[2][j] == symbol {
			return true
		}
	}

	if g.Board[0][0] == symbol && g.Board[1][1] == symbol && g.Board[2][2] == symbol {
		return true
	}
	if g.Board[0][2] == symbol && g.Board[1][1] == symbol && g.Board[2][0] == symbol {
		return true
	}

	return false
}

func (g *Game) isBoardFull() bool {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if g.Board[i][j] == 0 {
				return false
			}
		}
	}
	return true
}

func (g *Game) UpdateNextPlayer() {
	if g.Status == StatusPlayer1Turn {
		g.NextTurnID = g.Player1ID
	} else if g.Status == StatusPlayer2Turn {
		g.NextTurnID = g.Player2ID
	}
}

func (g *Game) SwitchPlayer() {
	if g.Status == StatusPlayer1Turn {
		g.Status = StatusPlayer2Turn
	} else if g.Status == StatusPlayer2Turn {
		g.Status = StatusPlayer1Turn
	}
}

func (g *Game) IsFinished() bool {
	return g.Status == StatusPlayer1Won || g.Status == StatusPlayer2Won || g.Status == StatusDraw
}

func (g *Game) GetWinnerID() string {
	switch g.Status {
	case StatusPlayer1Won:
		return g.Player1ID
	case StatusPlayer2Won:
		return g.Player2ID
	default:
		return ""
	}
}

func (g *Game) CanJoin() bool {
	return g.Mode == ModePvP && g.Player2ID == "" && !g.IsFinished()
}

func (g *Game) IsPlayerInGame(playerID string) bool {
	return g.Player1ID == playerID || g.Player2ID == playerID
}

func (g *Game) JoinGame(playerID string) error {
	if !g.CanJoin() {
		return ErrGameNotJoinable
	}
	if g.Player1ID == playerID {
		return ErrPlayerAlreadyInGame
	}

	g.Player2ID = playerID
	g.Status = StatusPlayer1Turn
	g.NextTurnID = g.Player1ID
	g.UpdatedAt = time.Now()

	return nil
}
