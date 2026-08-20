package dto

import "io"

type GetProfileRequest struct {
	UserId          string
	RequesterUserId string
}

type GetProfileResponse struct {
	Username     string `json:"username"`
	Nickname     string `json:"nickname"`
	Bio          string `json:"bio"`
	ProfileImage string `json:"profile_image"`
	IsOwner      bool   `json:"is_owner"`
	IsFollowing  bool   `json:"is_following"`
	Followers    int64  `json:"followers"`
	Following    int64  `json:"following"`
}

type GetMeRequest struct {
	UserID string `json:"-"`
}

type GetMeResponse struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	ProfileImage string `json:"profile_image"`
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

type UploadProfileImageRequest struct {
	UserID      string    `json:"-"`
	File        io.Reader `json:"-"`
	Filename    string    `json:"-"`
	ContentType string    `json:"-"`
	Size        int64     `json:"-"`
}

type UploadProfileImageResponse struct {
	ProfileImage string `json:"profile_image"`
}
