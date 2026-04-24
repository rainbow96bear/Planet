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
	TaskID uint `json:"-" `
	UserID uint `json:"-" `
}

type DeleteTaskResponse struct {
}
