package dto

import (
	"game/internal/domain/entities"
	"time"
)

type CreateGameRequest struct {
	Mode string `json:"mode" validate:"required,oneof=pvp pve"`
}

type MoveRequest struct {
	Row int `json:"row" validate:"required,min=0,max=2"`
	Col int `json:"col" validate:"required,min=0,max=2"`
}

type JoinGameRequest struct {
	GameID string `json:"game_id" validate:"required,uuid"`
}

type GameResponse struct {
	ID         string      `json:"id"`
	Board      [3][3]int   `json:"board"`
	Status     string      `json:"status"`
	Player1    PlayerInfo  `json:"player1"`
	Player2    *PlayerInfo `json:"player2,omitempty"`
	NextPlayer string      `json:"next_player"`
	Mode       string      `json:"mode"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
	Result     *GameResult `json:"result,omitempty"`
}

type PlayerInfo struct {
	ID    string `json:"id"`
	Login string `json:"login"`
	Score int    `json:"score"`
}

type GameResult struct {
	Winner      *PlayerInfo `json:"winner,omitempty"`
	Loser       *PlayerInfo `json:"loser,omitempty"`
	IsDraw      bool        `json:"is_draw"`
	WinningLine [][2]int    `json:"winning_line,omitempty"`
}

type GameShortInfo struct {
	ID      string      `json:"id"`
	Mode    string      `json:"mode"`
	Player1 PlayerInfo  `json:"player1"`
	Player2 *PlayerInfo `json:"player2,omitempty"`
	Status  string      `json:"status"`
}

type AvailableGamesResponse struct {
	Games []GameShortInfo `json:"games"`
	Total int             `json:"total"`
}

func ToGameResponse(game *entities.Game, player1, player2 *entities.User) GameResponse {
	resp := GameResponse{
		ID:         game.ID,
		Board:      game.Board,
		Status:     string(game.Status),
		NextPlayer: game.NextTurnID,
		Mode:       string(game.Mode),
		CreatedAt:  game.CreatedAt,
		UpdatedAt:  game.UpdatedAt,
		Player1: PlayerInfo{
			ID:    player1.UUID,
			Login: player1.Login,
			Score: player1.Score,
		},
	}

	if player2 != nil {
		resp.Player2 = &PlayerInfo{
			ID:    player2.UUID,
			Login: player2.Login,
			Score: player2.Score,
		}
	}

	if game.Status == entities.StatusPlayer1Won || game.Status == entities.StatusPlayer2Won || game.Status == entities.StatusDraw {
		result := GameResult{
			IsDraw: game.Status == entities.StatusDraw,
		}

		if !result.IsDraw {
			winner := player1
			loser := player2
			if game.Status == entities.StatusPlayer2Won && player2 != nil {
				winner = player2
				loser = player1
			}

			result.Winner = &PlayerInfo{
				ID:    winner.UUID,
				Login: winner.Login,
				Score: winner.Score,
			}

			if loser != nil {
				result.Loser = &PlayerInfo{
					ID:    loser.UUID,
					Login: loser.Login,
					Score: loser.Score,
				}
			}

			result.WinningLine = getWinningLine(game.Board)
		}

		resp.Result = &result
	}

	return resp
}

func ToGameShortInfo(game *entities.Game, player1, player2 *entities.User) GameShortInfo {
	info := GameShortInfo{
		ID:     game.ID,
		Mode:   string(game.Mode),
		Status: string(game.Status),
		Player1: PlayerInfo{
			ID:    player1.UUID,
			Login: player1.Login,
			Score: player1.Score,
		},
	}

	if player2 != nil {
		info.Player2 = &PlayerInfo{
			ID:    player2.UUID,
			Login: player2.Login,
			Score: player2.Score,
		}
	}

	return info
}

func getWinningLine(board [3][3]int) [][2]int {
	for i := 0; i < 3; i++ {
		if board[i][0] != 0 && board[i][0] == board[i][1] && board[i][1] == board[i][2] {
			return [][2]int{{i, 0}, {i, 1}, {i, 2}}
		}
	}

	for j := 0; j < 3; j++ {
		if board[0][j] != 0 && board[0][j] == board[1][j] && board[1][j] == board[2][j] {
			return [][2]int{{0, j}, {1, j}, {2, j}}
		}
	}

	if board[0][0] != 0 && board[0][0] == board[1][1] && board[1][1] == board[2][2] {
		return [][2]int{{0, 0}, {1, 1}, {2, 2}}
	}

	if board[0][2] != 0 && board[0][2] == board[1][1] && board[1][1] == board[2][0] {
		return [][2]int{{0, 2}, {1, 1}, {2, 0}}
	}

	return nil
}
