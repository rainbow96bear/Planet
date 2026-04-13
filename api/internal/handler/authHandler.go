package handler

import (
	"planet/internal/dto"
	"planet/internal/pkg"
	"planet/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler interface {
	CreateUser(c *gin.Context)
	CreateOAuthUser(c *gin.Context)
	CheckUsername(c *gin.Context)
	Login(c *gin.Context)
	OauthLogin(c *gin.Context)
}

type authHandler struct {
	authSvc service.AuthService
}

func NewAuthHandler(authSvc service.AuthService) AuthHandler {
	return &authHandler{authSvc: authSvc}
}

func (h *authHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Fail(c, 400, err.Error())
		return
	}

	user, err := h.authSvc.CreateUser(&req)
	if err != nil {
		pkg.Fail(c, 500, err.Error())
		return
	}

	pkg.Success(c, 201, user)
}

func (h *authHandler) CreateOAuthUser(c *gin.Context) {
	var req dto.CreateOAuthUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Fail(c, 400, err.Error())
		return
	}

	user, err := h.authSvc.CreateOAuthUser(&req)
	if err != nil {
		pkg.Fail(c, 401, err.Error())
		return
	}

	pkg.Success(c, 201, user)
}

func (h *authHandler) CheckUsername(c *gin.Context) {
	var req dto.CheckUsernameRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		pkg.Fail(c, 400, err.Error())
		return
	}

	res, err := h.authSvc.IsUsernameAvailable(&req)
	if err != nil {
		pkg.Fail(c, 500, err.Error())
		return
	}

	pkg.Success(c, 200, res)
}

func (h *authHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	res, err := h.authSvc.Login(&req)
	if err != nil {
		c.JSON(401, gin.H{"error": "invalid credentials"})
		return
	}

	c.JSON(200, res)
}

func (h *authHandler) OauthLogin(c *gin.Context) {
	var req dto.OauthLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	res, err := h.authSvc.OauthLogin(&req)
	if err != nil {
		c.JSON(401, gin.H{"error": ""})
	}

	c.JSON(200, res)
}
