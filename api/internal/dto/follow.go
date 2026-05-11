package dto

type FollowRequest struct {
	FollowerID  string
	FollowingID string
}

type FollowResponse struct {
	IsFollowing bool `json:"is_following"`
}

type UnfollowRequest struct {
	FollowerID  string
	FollowingID string
}

type UnfollowResponse struct {
	IsFollowing bool `json:"is_following"`
}

type IsFollowingRequest struct {
	FollowerID  string
	FollowingID string
}

type IsFollowingResponse struct {
	IsFollowing bool `json:"is_following"`
}
