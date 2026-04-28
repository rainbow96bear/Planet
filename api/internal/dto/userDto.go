package dto

type GetProfileRequest struct {
	UserId          uint
	RequesterUserId uint
}

type GetProfileResponse struct {
	Username  string `json:"username"`
	Nickname  string `json:"nickname"`
	Bio       string `json:"bio"`
	AvatarURL string `json:"avatar_url"`
	IsOwner   bool   `json:"is_owner"`
}

type GetMeRequest struct {
}

type GetMeResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
}

type GetProfileByUserIdRequest struct {
	UserId          uint
	RequesterUserId uint
}

type GetProfilByUserIdeResponse struct {
	Username  string `json:"username"`
	Nickname  string `json:"nickname"`
	Bio       string `json:"bio"`
	AvatarURL string `json:"avatar_url"`
	IsOwner   bool   `json:"is_owner"`
}
