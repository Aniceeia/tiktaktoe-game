// internal/interfaces/http/dto/game_dto.go
package dto

import (
	"game/internal/domain/entities"
	"time"
)

// CreateGameRequest - запрос на создание новой игры
type CreateGameRequest struct {
	Mode string `json:"mode" validate:"required,oneof=pvp pve"` // Режим игры: pvp (игрок vs игрок) или pve (игрок vs AI)
}

// MoveRequest - запрос на выполнение хода
type MoveRequest struct {
	Row int `json:"row" validate:"required,min=0,max=2"` // Строка (0-2)
	Col int `json:"col" validate:"required,min=0,max=2"` // Колонка (0-2)
}

// JoinGameRequest - запрос на присоединение к игре
type JoinGameRequest struct {
	GameID string `json:"game_id" validate:"required,uuid"` // ID игры для присоединения
}

// GameResponse - ответ с состоянием игры
type GameResponse struct {
	ID         string      `json:"id"`                // ID игры
	Board      [3][3]int   `json:"board"`             // Игровое поле (0 - пусто, 1 - X, 2 - O)
	Status     string      `json:"status"`            // Текущий статус игры
	Player1    PlayerInfo  `json:"player1"`           // Информация о первом игроке
	Player2    *PlayerInfo `json:"player2,omitempty"` // Информация о втором игроке (если есть)
	NextPlayer string      `json:"next_player"`       // ID игрока, чей ход следующий
	Mode       string      `json:"mode"`              // Режим игры (pvp/pve)
	CreatedAt  time.Time   `json:"created_at"`        // Время создания игры
	UpdatedAt  time.Time   `json:"updated_at"`        // Время последнего обновления
	Result     *GameResult `json:"result,omitempty"`  // Результат игры (если игра завершена)
}

// PlayerInfo - информация об игроке
type PlayerInfo struct {
	ID       string `json:"id"`       // ID игрока
	Username string `json:"username"` // Имя пользователя
	Score    int    `json:"score"`    // Счет игрока
}

// GameResult - результат завершенной игры
type GameResult struct {
	Winner      *PlayerInfo `json:"winner,omitempty"`       // Победитель (если есть)
	Loser       *PlayerInfo `json:"loser,omitempty"`        // Проигравший (если есть)
	IsDraw      bool        `json:"is_draw"`                // Ничья
	WinningLine [][2]int    `json:"winning_line,omitempty"` // Выигрышная линия (координаты клеток)
}

// GameShortInfo - краткая информация об игре (для списка игр)
type GameShortInfo struct {
	ID      string      `json:"id"`                // ID игры
	Mode    string      `json:"mode"`              // Режим игры
	Player1 PlayerInfo  `json:"player1"`           // Первый игрок
	Player2 *PlayerInfo `json:"player2,omitempty"` // Второй игрок (если есть)
	Status  string      `json:"status"`            // Статус игры
}

// AvailableGamesResponse - ответ со списком доступных игр
type AvailableGamesResponse struct {
	Games []GameShortInfo `json:"games"` // Список доступных игр
	Total int             `json:"total"` // Общее количество игр
}

// ErrorResponse - ответ с ошибкой
type ErrorResponse struct {
	Error   string `json:"error"`   // Сообщение об ошибке
	Code    string `json:"code"`    // Код ошибки
	Details string `json:"details"` // Детали ошибки (опционально)
}

// ToGameResponse - преобразует сущность Game в GameResponse
func ToGameResponse(game *entities.Game, player1, player2 *entities.User) GameResponse {
	resp := GameResponse{
		ID:         game.ID,
		Board:      game.Board,
		Status:     string(game.Status),
		NextPlayer: game.NextTurnID,
		Mode:       string(game.Mode),
		Player1: PlayerInfo{
			ID:       player1.UUID,
			Username: player1.Username,
			Score:    player1.Score,
		},
	}

	if player2 != nil {
		resp.Player2 = &PlayerInfo{
			ID:       player2.UUID,
			Username: player2.Username,
			Score:    player2.Score,
		}
	}

	// Если игра завершена, добавляем результат
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
				ID:       winner.UUID,
				Username: winner.Username,
				Score:    winner.Score,
			}

			if loser != nil {
				result.Loser = &PlayerInfo{
					ID:       loser.UUID,
					Username: loser.Username,
					Score:    loser.Score,
				}
			}

			result.WinningLine = getWinningLine(game.Board)
		}

		resp.Result = &result
	}

	return resp
}

// ToGameShortInfo - преобразует сущность Game в GameShortInfo
func ToGameShortInfo(game *entities.Game, player1, player2 *entities.User) GameShortInfo {
	info := GameShortInfo{
		ID:     game.ID,
		Mode:   string(game.Mode),
		Status: string(game.Status),
		Player1: PlayerInfo{
			ID:       player1.UUID,
			Username: player1.Username,
			Score:    player1.Score,
		},
	}

	if player2 != nil {
		info.Player2 = &PlayerInfo{
			ID:       player2.UUID,
			Username: player2.Username,
			Score:    player2.Score,
		}
	}

	return info
}

// getWinningLine - возвращает координаты выигрышной линии
func getWinningLine(board [3][3]int) [][2]int {
	// Проверка строк
	for i := 0; i < 3; i++ {
		if board[i][0] != 0 && board[i][0] == board[i][1] && board[i][1] == board[i][2] {
			return [][2]int{{i, 0}, {i, 1}, {i, 2}}
		}
	}

	// Проверка столбцов
	for j := 0; j < 3; j++ {
		if board[0][j] != 0 && board[0][j] == board[1][j] && board[1][j] == board[2][j] {
			return [][2]int{{0, j}, {1, j}, {2, j}}
		}
	}

	// Проверка диагоналей
	if board[0][0] != 0 && board[0][0] == board[1][1] && board[1][1] == board[2][2] {
		return [][2]int{{0, 0}, {1, 1}, {2, 2}}
	}

	if board[0][2] != 0 && board[0][2] == board[1][1] && board[1][1] == board[2][0] {
		return [][2]int{{0, 2}, {1, 1}, {2, 0}}
	}

	return nil
}
