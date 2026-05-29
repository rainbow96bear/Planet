package dto

import "time"

type GetFeedRequest struct {
}

type GetFeedResponse struct {
	FeedID        string    `json:"feed_id"`
	ActorID       string    `json:"actor_id"`
	ActorNickname string    `json:"actor_nickname"`
	Type          string    `json:"type"`
	TargetID      string    `json:"target_id"`
	TargetType    string    `json:"target_type"`
	TaskTitle     *string   `json:"task_title,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}
