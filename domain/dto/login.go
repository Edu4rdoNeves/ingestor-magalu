package dto

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginAuth struct {
	Token     string
	ExpiresAt int64
}
