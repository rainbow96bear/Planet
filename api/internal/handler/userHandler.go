package handler

import (
	"planet/internal/dto"
	"planet/internal/pkg"
	"planet/internal/service"

	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	CreateUser(c *gin.Context)
	CheckUsername(c *gin.Context)
	Login(c *gin.Context)
}

type userHandler struct {
	userSvc service.UserService
}

func NewUserHandler(userSvc service.UserService) UserHandler {
	return &userHandler{userSvc: userSvc}
}

func (h *userHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Fail(c, 400, err.Error())
		return
	}

	user, err := h.userSvc.CreateUser(&req)
	if err != nil {
		pkg.Fail(c, 500, err.Error())
		return
	}

	pkg.Success(c, 201, user)
}

func (h *userHandler) CheckUsername(c *gin.Context) {
	var req dto.CheckUsernameRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		pkg.Fail(c, 400, err.Error())
		return
	}

	res, err := h.userSvc.IsUsernameAvailable(&req)
	if err != nil {
		pkg.Fail(c, 500, err.Error())
		return
	}

	pkg.Success(c, 200, res)
}

func (h *userHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	res, err := h.userSvc.Login(&req)
	if err != nil {
		c.JSON(401, gin.H{"error": "invalid credentials"})
		return
	}

	// httpOnly Cookie에 토큰 저장
	c.SetCookie("access_token", res.AccessToken, 60*60, "/", "", true, true)
	c.SetCookie("refresh_token", res.RefreshToken, 60*60*24*7, "/", "", true, true)
	//                                         ↑ 만료시간(초)       ↑secure ↑httpOnly

	c.JSON(200, dto.LoginResponse{
		Username: req.Username,
	})
}
