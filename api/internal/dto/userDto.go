package dto

type GetProfileRequest struct {
	UserId          string
	RequesterUserId string
}

type GetProfileResponse struct {
	Username    string `json:"username"`
	Nickname    string `json:"nickname"`
	IsOwner     bool   `json:"is_owner"`
	IsFollowing bool   `json:"is_following"`
	Followers   int64  `json:"followers"`
	Following   int64  `json:"following"`
}

type GetMeRequest struct {
}

type GetMeResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type GetProfileByUserIdRequest struct {
	UserId          string
	RequesterUserId string
}

type GetProfilByUserIdeResponse struct {
	Username  string `json:"username"`
	Nickname  string `json:"nickname"`
	Bio       string `json:"bio"`
	AvatarURL string `json:"avatar_url"`
	IsOwner   bool   `json:"is_owner"`
}

type UpdateProfileRequest struct {
	UserID   string `json:"-"`
	Nickname string `json:"nickname" binding:"required"`
}

type UpdateProfileResponse struct {
	UserID   string `json:"userid"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
}
