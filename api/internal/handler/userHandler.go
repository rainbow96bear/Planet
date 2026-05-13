package handler

import (
	"fmt"
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
	GetActivity(c *gin.Context)
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
	fmt.Printf("user id : %+v, requester user id : %+v\n", req.UserId, req.RequesterUserId)
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

func (h *userHandler) GetActivity(c *gin.Context) {
	userID := c.Param("userid")
	requesterUserId := c.GetString("userID")

	req := &dto.GetActivityRequest{
		UserID:          userID,
		RequesterUserId: requesterUserId,
	}

	result, err := h.userSvc.GetActivity(req)
	if err != nil {
		pkg.Fail(c, 500, err.Error())
		return
	}

	pkg.Success(c, 200, result)
}
