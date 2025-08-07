package api

type CreateGameRequest struct {
	Mode string `json:"mode" validate:"required,oneof=pvp pve"`
}

type UpdateGameRequest struct {
	Board [3][3]int `json:"board" validate:"required,row=3,col=3"`
}

type GameResponse struct {
	ID         string    `json:"id"`
	Board      [3][3]int `json:"board"`
	Status     string    `json:"status"`
	Player1    string    `json:"player1"`
	Player2    string    `json:"player2"`
	NextPlayer string    `json:"next_player"`
	Mode       string    `json:"mode"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type User struct {
	UUID         string `json:"uuid"`
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
	Score        int    `json:"score"`
}

type UserResponse struct {
	UUID     string `json:"uuid"`
	Username string `json:"username"`
	Score    int    `json:"score"`
}

type SignUpRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// last
type GameShortInfo struct {
	ID      string `json:"id"`
	Player1 string `json:"player1"`
	Mode    string `json:"mode"`
}
