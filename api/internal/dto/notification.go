package dto

import "time"

type GetNotificationsResponse struct {
	ID            string    `json:"id"`
	ActorID       string    `json:"actor_id"`
	ActorNickname string    `json:"actor_nickname"`
	Type          string    `json:"type"`
	IsRead        bool      `json:"is_read"`
	CreatedAt     time.Time `json:"created_at"`
}

type GetUnreadCountResponse struct {
	Count int64 `json:"count"`
}
