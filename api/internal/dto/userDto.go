package dto

type GetProfileRequest struct {
	Username          string
	RequesterUsername string
}

type GetProfileResponse struct {
	Username  string `json:"username"`
	Nickname  string `json:"nickname"`
	Bio       string `json:"bio"`
	AvatarURL string `json:"avatar_url"`
	IsOwner   bool   `json:"is_owner"`
}
