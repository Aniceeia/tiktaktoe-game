package dto

type SignUpRequest struct {
	Login    string `json:"login" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=6"`
}

type UserResponse struct {
	UUID  string `json:"uuid"`
	Login string `json:"login"`
	Score int    `json:"score"`
}

type LoginResponse struct {
	UUID string `json:"uuid"`
}
