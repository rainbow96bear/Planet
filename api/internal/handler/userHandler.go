package handler

import (
	"planet/internal/dto"
	"planet/internal/pkg"
	"planet/internal/service"

	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	GetProfile(c *gin.Context)
	Follow(c *gin.Context)
	Unfollow(c *gin.Context)
	UpdateProfile(c *gin.Context)
	UploadProfileImage(c *gin.Context)
	DeleteProfileImage(c *gin.Context)
}

type userHandler struct {
	userSvc service.UserService
}

func NewUserHandler(userSvc service.UserService) UserHandler {
	return &userHandler{
		userSvc: userSvc,
	}
}

func (h *userHandler) GetProfile(c *gin.Context) {
	userID := c.Param("userid")

	req := dto.GetProfileRequest{
		UserId:          userID,
		RequesterUserId: c.GetString("userID"),
	}
	profile, err := h.userSvc.GetProfile(&req)
	if err != nil {
		pkg.Fail(c, 500, err.Error())
		return
	}

	pkg.Success(c, 200, profile)
}

func (h *userHandler) Follow(c *gin.Context) {
	followingID := c.Param("userid")
	followerID := c.GetString("userID") // JWT 미들웨어에서 세팅한 값

	req := dto.FollowRequest{
		FollowerID:  followerID,
		FollowingID: followingID,
	}

	result, err := h.userSvc.Follow(&req)
	if err != nil {
		pkg.Fail(c, 409, err.Error())
		return
	}
	pkg.Success(c, 201, result)
}

func (h *userHandler) Unfollow(c *gin.Context) {
	followingID := c.Param("userid")

	followerID := c.GetString("userID") // JWT 미들웨어에서 세팅한 값

	req := dto.UnfollowRequest{
		FollowerID:  followerID,
		FollowingID: followingID,
	}

	result, err := h.userSvc.Unfollow(&req)
	if err != nil {
		pkg.Fail(c, 400, err.Error())
		return
	}
	pkg.Success(c, 200, result)
}

func (h *userHandler) UpdateProfile(c *gin.Context) {
	userid := c.Param("userid")
	requesterID := c.GetString("userID")

	if userid != requesterID {
		pkg.Fail(c, 403, "권한이 없습니다")
		return
	}

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Fail(c, 400, err.Error())
		return
	}
	req.UserID = userid

	result, err := h.userSvc.UpdateProfile(&req)
	if err != nil {
		pkg.Fail(c, 500, err.Error())
		return
	}
	pkg.Success(c, 200, result)
}

func (h *userHandler) UploadProfileImage(c *gin.Context) {
	userid := c.Param("userid")
	requesterID := c.GetString("userID")

	if userid != requesterID {
		pkg.Fail(c, 403, "권한이 없습니다")
		return
	}

	fileHeader, err := c.FormFile("profile_image")
	if err != nil {
		pkg.Fail(c, 400, "profile_image 파일이 필요합니다")
		return
	}

	// 크기/타입은 클라이언트에서도 검증하지만, 서버에서 재검증 (우회 가능하므로 필수)
	const maxProfileImageSize = 5 * 1024 * 1024 // 5MB
	if fileHeader.Size > maxProfileImageSize {
		pkg.Fail(c, 400, "이미지 크기는 5MB 이하여야 합니다")
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		pkg.Fail(c, 400, "파일을 열 수 없습니다")
		return
	}
	defer file.Close()

	contentType := fileHeader.Header.Get("Content-Type")
	if !isAllowedImageType(contentType) {
		pkg.Fail(c, 400, "지원하지 않는 이미지 형식입니다")
		return
	}

	req := dto.UploadProfileImageRequest{
		UserID:      userid,
		File:        file,
		Filename:    fileHeader.Filename,
		ContentType: contentType,
		Size:        fileHeader.Size,
	}

	result, err := h.userSvc.UploadProfileImage(&req)
	if err != nil {
		pkg.Fail(c, 500, "이미지 업로드에 실패했습니다")
		return
	}

	pkg.Success(c, 200, result)
}

func (h *userHandler) DeleteProfileImage(c *gin.Context) {
	userid := c.Param("userid")
	requesterID := c.GetString("userID")

	if userid != requesterID {
		pkg.Fail(c, 403, "권한이 없습니다")
		return
	}

	if err := h.userSvc.DeleteProfileImage(userid); err != nil {
		pkg.Fail(c, 500, "이미지 삭제에 실패했습니다")
		return
	}

	pkg.Success(c, 200, nil)
}

func isAllowedImageType(contentType string) bool {
	switch contentType {
	case "image/jpeg", "image/png", "image/webp", "image/gif":
		return true
	default:
		return false
	}
}
