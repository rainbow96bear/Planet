package dto

type SearchUsersRequest struct {
	Q               string `form:"q" binding:"required"`
	RequesterUserId uint
}

type SearchUsersResponse struct {
	UserId      uint   `json:"userid"`
	Username    string `json:"username"`
	Nickname    string `json:"nickname"`
	IsFollowing bool   `json:"is_following"`
}
