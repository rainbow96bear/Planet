package dto

import "time"

type CreateTaskRequest struct {
	UserID      string    `json:"-" `
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	Date        time.Time `json:"date" binding:"required"`
	IsPublic    bool      `json:"is_public"`
}

type CreateTaskResponse struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	IsCompleted bool      `json:"is_completed"`
	IsPublic    bool      `json:"is_public"`
}

type DeleteTaskRequest struct {
	ID     string `json:"-" `
	UserID string `json:"-" `
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
	Date        time.Time `json:"date"`
	IsCompleted bool      `json:"is_completed"`
	IsPublic    bool      `json:"is_public"`
}

type ToggleTaskRequest struct {
	ID string `json:"id"`
}

type ToggleTaskResponse struct {
	ID          string `json:"id"`
	IsCompleted bool   `json:"is_completed"`
}
