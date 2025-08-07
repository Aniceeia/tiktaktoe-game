package api

import "game/internal/service"

func ToResponse(game *service.Game) GameResponse {
	return GameResponse{
		ID:         game.ID,
		Board:      game.Board,
		Status:     game.Status,
		Player1:    game.Player1,
		Player2:    game.Player2,
		NextPlayer: game.NextTurn,
		Mode:       game.Mode,
	}
}

// last
func ToShortResponse(game *service.Game) GameShortInfo {
	return GameShortInfo{
		ID:      game.ID,
		Player1: game.Player1,
		Mode:    game.Mode,
	}
}
