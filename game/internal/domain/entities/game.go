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
	ID         string     `json:"id"`
	Board      [3][3]int  `json:"board"`
	Status     GameStatus `json:"status"`
	Player1ID  string     `json:"player1_id"`
	Player2ID  string     `json:"player2_id"`
	NextTurnID string     `json:"next_turn_id"`
	Mode       GameMode   `json:"mode"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// 1 - player "X" (Player1)
// 2 - player "O" (Player2/AI)

func NewGame(creator string, mode GameMode) *Game {
	now := time.Now()
	return &Game{
		ID:         generateID(),
		Board:      [3][3]int{},
		Status:     StatusPlayer1Turn,
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
	g.switchPlayer()
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
	// Проверка строк
	for i := 0; i < 3; i++ {
		if g.Board[i][0] == symbol && g.Board[i][1] == symbol && g.Board[i][2] == symbol {
			return true
		}
	}

	// Проверка столбцов
	for j := 0; j < 3; j++ {
		if g.Board[0][j] == symbol && g.Board[1][j] == symbol && g.Board[2][j] == symbol {
			return true
		}
	}

	// Проверка диагоналей
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

func (g *Game) switchPlayer() {
	if g.Status == StatusPlayer1Turn {
		g.Status = StatusPlayer2Turn
	} else if g.Status == StatusPlayer2Turn {
		g.Status = StatusPlayer1Turn
	}
}

// IsFinished возвращает true, если игра завершена
func (g *Game) IsFinished() bool {
	return g.Status == StatusPlayer1Won || g.Status == StatusPlayer2Won || g.Status == StatusDraw
}

// GetWinnerID возвращает ID победителя или пустую строку
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

// CanJoin возвращает true, если к игре можно присоединиться
func (g *Game) CanJoin() bool {
	return g.Mode == ModePvP && g.Player2ID == "" && !g.IsFinished()
}

// IsPlayerInGame проверяет, является ли пользователь участником игры
func (g *Game) IsPlayerInGame(playerID string) bool {
	return g.Player1ID == playerID || g.Player2ID == playerID
}
