package dto

import "time"

type GetFeedRequest struct {
}

type GetFeedResponse struct {
	FeedID        string    `json:"feed_id"`
	ActorID       string    `json:"actor_id"`
	ActorNickname string    `json:"actor_nickname"`
	Type          string    `json:"type"`
	TaskID        string    `json:"task_id"`
	TaskTitle     *string   `json:"task_title,omitempty"`
	LikeCount     int       `json:"like_count"`
	CheerCount    int       `json:"cheer_count"`
	IsLiked       bool      `json:"is_liked"`
	IsCheered     bool      `json:"is_cheered"`
	CreatedAt     time.Time `json:"created_at"`
}
