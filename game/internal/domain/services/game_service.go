package services

import (
	"context"
	"game/internal/domain/entities"
	"game/internal/domain/repositories"
)

type GameService struct {
	gameRepo  repositories.GameRepository
	userRepo  repositories.UserRepository
	aiService *AIService
}

func NewGameService(gameRepo repositories.GameRepository, userRepo repositories.UserRepository, aiService *AIService) *GameService {
	return &GameService{
		gameRepo:  gameRepo,
		userRepo:  userRepo,
		aiService: aiService,
	}
}

func (s *GameService) CreateGame(ctx context.Context, playerID string, mode entities.GameMode) (*entities.Game, error) {
	// Проверяем, что пользователь существует
	_, err := s.userRepo.GetByID(ctx, playerID)
	if err != nil {
		return nil, entities.ErrUserNotFound
	}

	game := entities.NewGame(playerID, mode)

	if mode == entities.ModePvE {
		game.Player2ID = "AI"
	}

	err = s.gameRepo.Create(ctx, game)
	if err != nil {
		return nil, err
	}

	return game, nil
}

func (s *GameService) JoinGame(ctx context.Context, gameID, playerID string) error {
	// Проверяем, что пользователь существует
	_, err := s.userRepo.GetByID(ctx, playerID)
	if err != nil {
		return entities.ErrUserNotFound
	}

	game, err := s.gameRepo.GetByID(ctx, gameID)
	if err != nil {
		return entities.ErrGameNotFound
	}

	if !game.CanJoin() {
		return entities.ErrGameFull
	}

	game.Player2ID = playerID
	game.Status = entities.StatusPlayer2Turn
	game.NextTurnID = playerID
	game.UpdatedAt = entities.GetCurrentTime()

	return s.gameRepo.Update(ctx, game)
}

func (s *GameService) MakeMove(ctx context.Context, gameID, playerID string, row, col int) (*entities.Game, error) {
	game, err := s.gameRepo.GetByID(ctx, gameID)
	if err != nil {
		return nil, entities.ErrGameNotFound
	}

	// Проверяем, что пользователь является участником игры
	if !game.IsPlayerInGame(playerID) {
		return nil, entities.ErrForbidden
	}

	err = game.MakeMove(playerID, row, col)
	if err != nil {
		return nil, err
	}

	// Если игра в PvE режиме и ход сделал игрок, делаем ход AI
	if game.Mode == entities.ModePvE && game.Player1ID == playerID && game.Status == entities.StatusPlayer2Turn && !game.IsFinished() {
		s.makeAIMove(game)
	}

	err = s.gameRepo.Update(ctx, game)
	if err != nil {
		return nil, err
	}

	// Обновляем счет победителя
	if game.IsFinished() {
		winnerID := game.GetWinnerID()
		if winnerID != "" && winnerID != "AI" {
			err = s.userRepo.IncreaseScore(ctx, winnerID)
			if err != nil {
				// Логируем ошибку, но не прерываем выполнение
				// TODO: добавить proper logging
			}
		}
	}

	return game, nil
}

func (s *GameService) GetGameState(ctx context.Context, gameID string) (*entities.Game, error) {
	return s.gameRepo.GetByID(ctx, gameID)
}

func (s *GameService) GetAvailableGames(ctx context.Context) ([]*entities.Game, error) {
	return s.gameRepo.GetAvailableGames(ctx)
}

func (s *GameService) GetUserGames(ctx context.Context, userID string) ([]*entities.Game, error) {
	return s.gameRepo.GetByPlayerID(ctx, userID)
}

func (s *GameService) makeAIMove(game *entities.Game) {
	if game.Status != entities.StatusPlayer2Turn || game.IsFinished() {
		return
	}

	row, col := s.aiService.FindBestMove(game.Board)
	game.Board[row][col] = 2 // O для AI
	game.UpdateStatus()
	game.UpdateNextPlayer()
	game.UpdatedAt = entities.GetCurrentTime()
}
