package dto

import "time"

type CreateTaskRequest struct {
	UserID      uint      `json:"-" `
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	Date        time.Time `json:"date" binding:"required"`
	IsPublic    bool      `json:"is_public"`
}

type CreateTaskResponse struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Date        time.Time `json:"date"`
	IsCompleted bool      `json:"is_completed"`
	IsPublic    bool      `json:"is_public"`
}

type DeleteTaskRequest struct {
	ID     uint `json:"-" `
	UserID uint `json:"-" `
}

type DeleteTaskResponse struct {
}

type GetTasksByMonthRequest struct {
	UserID          uint `form:"-"`
	RequesterUserId uint `form:"-"`
	Year            int  `form:"year" binding:"required"`
	Month           int  `form:"month" binding:"required,min=1,max=12"`
}

type GetTasksByMonthResponse struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	IsCompleted bool      `json:"is_completed"`
	IsPublic    bool      `json:"is_public"`
}

type ToggleTaskRequest struct {
	ID uint `json:"id"`
}

type ToggleTaskResponse struct {
	ID          uint `json:"id"`
	IsCompleted bool `json:"is_completed"`
}
