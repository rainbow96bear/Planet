package dto

type SearchUsersRequest struct {
	Q               string `form:"q" binding:"required"`
	RequesterUserId string
}

type SearchUsersResponse struct {
	UserId      string `json:"userid"`
	Username    string `json:"username"`
	Nickname    string `json:"nickname"`
	IsFollowing bool   `json:"is_following"`
}
