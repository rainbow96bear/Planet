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
