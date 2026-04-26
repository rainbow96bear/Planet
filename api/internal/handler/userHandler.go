package handler

import (
	"planet/internal/dto"
	"planet/internal/pkg"
	"planet/internal/service"

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
	req := dto.GetProfileRequest{
		Username:          c.Param("username"),
		RequesterUsername: c.GetString("username"), // 비로그인이면 ""
	}

	profile, err := h.userSvc.GetProfile(&req)
	if err != nil {
		pkg.Fail(c, 500, err.Error())
		return
	}

	pkg.Success(c, 200, profile)
}
