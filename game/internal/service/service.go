package service

type GameService interface {
	CreateGame(creator string, mode string) (*Game, error)
	JoinGame(gameID string, joiner string) error
	ProcessMove(gameID string, user string, newBoard [3][3]int) (*Game, error)
	GetGameState(gameID string) (*Game, error)
	GetAvailableGames() ([]*Game, error)
}

type GameServiceImpl struct {
	repo GameRepository
	ai   *AI
}

func NewGameService(repo GameRepository, ai *AI) GameService {
	return &GameServiceImpl{repo, ai}
}

func (s *GameServiceImpl) CreateGame(creator string, mode string) (*Game, error) {
	game := &Game{
		ID:       GenerateID(),
		Board:    [3][3]int{},
		Status:   StatusWaiting,
		Player1:  creator,
		NextTurn: creator,
		Mode:     mode,
	}

	if mode == "pve" {
		game.Player2 = "AI"
	}

	s.repo.Save(game)
	return game, nil
}

func (s *GameServiceImpl) JoinGame(gameID string, joiner string) error {
	game, err := s.repo.Get(gameID)
	if err != nil {
		return errGameNotFound
	}

	if game.Mode != "pvp" {
		return errInvalidMode
	}

	if game.Player2 != "" {
		return errGameFull
	}

	game.Player2 = joiner
	s.repo.Save(game)
	return nil
}

func (s *GameServiceImpl) ProcessMove(gameID string, user string, newBoard [3][3]int) (*Game, error) {
	currentGame, err := s.repo.Get(gameID)
	if err != nil {
		return nil, errGameNotFound
	}

	if currentGame.Status == StatusPlayer1Won || currentGame.Status == StatusDraw || currentGame.Status == StatusPlayer2Won {
		return nil, errGameCompleted
	}

	// Check if it's the user's turn
	if currentGame.NextTurn != user {
		return nil, errWrongTurn
	}

	if err := s.isValidBoard(newBoard); err != nil {
		return nil, err
	}

	if !isValidMove(currentGame.Board, newBoard) {
		return nil, errInvalidMove
	}

	currentGame.Board = newBoard
	currentGame.updateGameStatus()

	// Handle turn switching
	if currentGame.Status == StatusWaiting {
		if currentGame.Mode == "pve" {
			// Switch to AI turn
			currentGame.NextTurn = "AI"
			s.makeComputerMove(currentGame)
			currentGame.updateGameStatus()
			currentGame.NextTurn = user // Switch back to player
			currentGame.Status = StatusPlayer1Turn
		} else {
			// Switch to other player
			if currentGame.NextTurn == currentGame.Player1 {
				currentGame.NextTurn = currentGame.Player2
				currentGame.Status = StatusPlayer2Turn
			} else {
				currentGame.NextTurn = currentGame.Player1
				currentGame.Status = StatusPlayer1Turn
			}
		}
	}

	s.repo.Save(currentGame)
	return currentGame, nil
}

func (s *GameServiceImpl) GetGameState(gameID string) (*Game, error) {
	return s.repo.Get(gameID)
}

func (s *GameServiceImpl) isValidBoard(board [3][3]int) error {
	if len(board) != 3 {
		return errInvalidBoard
	}
	for _, row := range board {
		if len(row) != 3 {
			return errInvalidBoard
		}
	}
	return nil
}

func isValidMove(oldBoard, newBoard [3][3]int) bool {
	diffCount := 0
	for i := range 3 {
		for j := range 3 {
			old := oldBoard[i][j]
			new := newBoard[i][j]
			if new != old {
				if old != empty && new != x {
					return false
				}
				if old != empty && new != o {
					return false
				}
				diffCount++
			}

		}
	}
	return diffCount == 1
}

func (s *GameServiceImpl) makeComputerMove(game *Game) {
	row, col := s.ai.FindBestMove(game.Board)
	game.Board[row][col] = o
}

// last
func (s *GameServiceImpl) GetAvailableGames() ([]*Game, error) {
	allGames, _ := s.repo.GetAllGames()
	var availableGames []*Game
	for _, game := range allGames {
		if game.Status == StatusWaiting && game.Player2 == "" {
			availableGames = append(availableGames, game)
		}
	}
	return availableGames, nil
}
