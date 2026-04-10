package dto

type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Nickname string `json:"nickname" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}
