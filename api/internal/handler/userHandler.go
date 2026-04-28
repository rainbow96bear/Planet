package handler

import (
	"net/http"
	"planet/internal/dto"
	"planet/internal/pkg"
	"planet/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	GetProfile(c *gin.Context)
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

	profile, err := h.userSvc.GetProfile(&req)
	if err != nil {
		pkg.Fail(c, 500, err.Error())
		return
	}

	pkg.Success(c, 200, profile)
}
