package service

// MockGameService реализует GameService для тестов
type MockGameService struct {
	CreateGameFunc        func(creator string, mode string) (*Game, error)
	JoinGameFunc          func(gameID string, joiner string) error
	ProcessMoveFunc       func(gameID string, user string, newBoard [3][3]int) (*Game, error)
	GetGameStateFunc      func(gameID string) (*Game, error)
	GetAvailableGamesFunc func() ([]*Game, error)
}

func (m *MockGameService) CreateGame(creator string, mode string) (*Game, error) {
	return m.CreateGameFunc(creator, mode)
}

func (m *MockGameService) JoinGame(gameID string, joiner string) error {
	return m.JoinGameFunc(gameID, joiner)
}

func (m *MockGameService) ProcessMove(gameID string, user string, newBoard [3][3]int) (*Game, error) {
	return m.ProcessMoveFunc(gameID, user, newBoard)
}

func (m *MockGameService) GetGameState(gameID string) (*Game, error) {
	return m.GetGameStateFunc(gameID)
}

func (m *MockGameService) GetAvailableGames() ([]*Game, error) {
	return m.GetAvailableGamesFunc()
}

// MockAuthService реализует AuthService для тестов
type MockAuthService struct {
	RegisterFunc      func(req SignUpRequest) error
	AuthenticateFunc  func(login, password string) (string, error)
	GetUserByUUIDFunc func(uuid string) (*User, error)
}

func (m *MockAuthService) Register(req SignUpRequest) error {
	return m.RegisterFunc(req)
}

func (m *MockAuthService) Authenticate(login, password string) (string, error) {
	return m.AuthenticateFunc(login, password)
}

func (m *MockAuthService) GetUserByUUID(uuid string) (*User, error) {
	return m.GetUserByUUIDFunc(uuid)
}
