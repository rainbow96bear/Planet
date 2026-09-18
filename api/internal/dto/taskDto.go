package dto

import "time"

type CreateTaskRequest struct {
	UserID      string    `json:"-"`
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	StartAt     time.Time `json:"start_at" binding:"required"`
	EndAt       time.Time `json:"end_at" binding:"required,gtfield=StartAt"`
	IsPublic    bool      `json:"is_public"`
}

type CreateTaskResponse struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartAt     time.Time `json:"start_at"`
	EndAt       time.Time `json:"end_at"`
	IsCompleted bool      `json:"is_completed"`
	IsPublic    bool      `json:"is_public"`
}

type UpdateTaskRequest struct {
	ID          string    `json:"-"`
	UserID      string    `json:"-"`
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	StartAt     time.Time `json:"start_at" binding:"required"`
	EndAt       time.Time `json:"end_at" binding:"required,gtfield=StartAt"`
	IsPublic    bool      `json:"is_public"`
}

type UpdateTaskResponse struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartAt     time.Time `json:"start_at"`
	EndAt       time.Time `json:"end_at"`
	IsCompleted bool      `json:"is_completed"`
	IsPublic    bool      `json:"is_public"`
}

type DeleteTaskRequest struct {
	ID     string `json:"-"`
	UserID string `json:"-"`
}

type DeleteTaskResponse struct {
}

type GetTasksByMonthRequest struct {
	UserID          string `form:"-"`
	RequesterUserId string `form:"-"`
	Year            int    `form:"year" binding:"required"`
	Month           int    `form:"month" binding:"required,min=1,max=12"`
}

type GetTasksByMonthResponse struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartAt     time.Time `json:"start_at"`
	EndAt       time.Time `json:"end_at"`
	IsCompleted bool      `json:"is_completed"`
	IsPublic    bool      `json:"is_public"`
}

type ToggleTaskRequest struct {
	ID     string `json:"-"`
	UserID string `json:"-"`
}

type ToggleTaskResponse struct {
	ID          string `json:"id"`
	IsCompleted bool   `json:"is_completed"`
}

type OrbitScheduleResponse struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	Nickname     string    `json:"nickname"`
	ProfileImage string    `json:"profile_image"`
	Title        string    `json:"title"`
	StartAt      time.Time `json:"start_at"`
	EndAt        time.Time `json:"end_at"`
}

type GetOrbitSchedulesByMonthRequest struct {
	OrbiterID string `form:"-"`
	Year      int    `form:"year" binding:"required"`
	Month     int    `form:"month" binding:"required,min=1,max=12"`
}
