package dto

import "planet/internal/model"

type AddReactionRequest struct {
	TaskID      string             `json:"-"`
	UserID      string             `json:"-"`
	TaskOwnerID string             `json:"-"`
	Type        model.ReactionType `json:"type" binding:"required,oneof=like cheer"`
}

type AddReactionResponse struct {
	TaskID string             `json:"task_id"`
	Type   model.ReactionType `json:"type"`
}

type RemoveReactionRequest struct {
	TaskID string `json:"-"`
	UserID string `json:"-"`
	Type   string `json:"type" binding:"required,oneof=like cheer"`
}
