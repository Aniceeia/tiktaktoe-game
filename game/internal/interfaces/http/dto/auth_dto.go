package dto

// SignUpRequest - запрос на регистрацию
type SignUpRequest struct {
	Login    string `json:"login" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=6"`
}

// LoginRequest - запрос на вход
type LoginRequest struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// UserResponse - ответ с информацией о пользователе
type UserResponse struct {
	UUID     string `json:"uuid"`
	Username string `json:"username"`
	Score    int    `json:"score"`
}

// LoginResponse - ответ на вход
type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}
