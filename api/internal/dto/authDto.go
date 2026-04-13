package dto

import "time"

type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=4"`
	Nickname string `json:"nickname" binding:"required,min=4"`
	Password string `json:"password" binding:"required,min=8"`
}

type CreateUserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Nickname  string    `json:"nickname"`
	CreatedAt time.Time `json:"createdAt"`
}

type CreateOAuthUserRequest struct {
	Username  string `json:"username"   binding:"required,min=4"`
	Nickname  string `json:"nickname"   binding:"required,min=4"`
	TempToken string `json:"temp_token" binding:"required"`
}

type CreateOAuthUserResponse struct {
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
	Nickname     string `json:"nickname"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type OauthLoginRequest struct {
	Provider   string `json:"provider" bindig:"required"`
	ProviderID string `json:"provider_id" binding:"required"`
}

type OauthLoginResponse struct {
	IsNewUser    bool   `json:"is_new_user"`
	TempToken    string `json:"temp_token,omitempty"`    // 신규 유저만
	Username     string `json:"username,omitempty"`      // 기존 유저만
	Nickname     string `json:"nickname,omitempty"`      // 기존 유저만
	AccessToken  string `json:"access_token,omitempty"`  // 기존 유저만
	RefreshToken string `json:"refresh_token,omitempty"` // 기존 유저만
}
