package handler

import (
	"fmt"
	"net/http"
	"planet/internal/dto"
	"planet/internal/pkg"
	"planet/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	GetProfile(c *gin.Context)
	Follow(c *gin.Context)
	Unfollow(c *gin.Context)
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
	userID, err := strconv.ParseUint(c.Param("userid"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	req := dto.GetProfileRequest{
		UserId:          uint(userID),
		RequesterUserId: c.GetUint("userID"),
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
	followingID, err := strconv.ParseUint(c.Param("userid"), 10, 64)
	if err != nil {
		pkg.Fail(c, 400, "잘못된 요청입니다")
		return
	}

	followerID := c.GetUint("userID") // JWT 미들웨어에서 세팅한 값

	req := dto.FollowRequest{
		FollowerID:  followerID,
		FollowingID: uint(followingID),
	}

	result, err := h.userSvc.Follow(&req)
	if err != nil {
		pkg.Fail(c, 409, err.Error())
		return
	}
	pkg.Success(c, 201, result)
}

func (h *userHandler) Unfollow(c *gin.Context) {
	followingID, err := strconv.ParseUint(c.Param("userid"), 10, 64)
	if err != nil {
		pkg.Fail(c, 400, "잘못된 요청입니다")
		return
	}

	followerID := c.GetUint("userID") // JWT 미들웨어에서 세팅한 값

	req := dto.UnfollowRequest{
		FollowerID:  followerID,
		FollowingID: uint(followingID),
	}

	result, err := h.userSvc.Unfollow(&req)
	if err != nil {
		pkg.Fail(c, 400, err.Error())
		return
	}
	pkg.Success(c, 200, result)
}
