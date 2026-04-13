package dto

import "time"

type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Nickname string `json:"nickname" binding:"required,min=2"`
	Password string `json:"password" binding:"required,min=8"`
}

type CreateUserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Nickname  string    `json:"nickname"`
	CreatedAt time.Time `json:"createdAt"`
}

type CheckUsernameRequest struct {
	Username string `form:"username" binding:"required"`
}

type CheckUsernameResponse struct {
	Username  string `json:"username"`
	Available bool   `json:"available"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Username     string `json:"username"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
