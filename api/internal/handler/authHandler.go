package handler

import (
	"planet/internal/dto"
	"planet/internal/pkg"
	"planet/internal/service"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthHandler interface {
	CreateUser(c *gin.Context)
	CreateOAuthUser(c *gin.Context)
	CheckUsername(c *gin.Context)
	Login(c *gin.Context)
	OauthLogin(c *gin.Context)
	Refresh(c *gin.Context)
}

type authHandler struct {
	authSvc service.AuthService
}

func NewAuthHandler(authSvc service.AuthService) AuthHandler {
	return &authHandler{authSvc: authSvc}
}

func (h *authHandler) CreateUser(c *gin.Context) {
	req := dto.CreateUserRequest{
		Username:     c.PostForm("username"),
		Nickname:     c.PostForm("nickname"),
		Password:     c.PostForm("password"),
		AgreeTerms:   c.PostForm("agreeTerms") == "true",
		AgreePrivacy: c.PostForm("agreePrivacy") == "true",
	}

	if len(req.Username) < 4 || len(req.Username) > 20 {
		pkg.Fail(c, 400, "username은 4~20자여야 합니다")
		return
	}
	if len(req.Nickname) < 2 || len(req.Nickname) > 20 {
		pkg.Fail(c, 400, "nickname은 2~20자여야 합니다")
		return
	}
	if len(req.Password) < 8 || len(req.Password) > 72 {
		pkg.Fail(c, 400, "password는 8~72자여야 합니다")
		return
	}

	// 이미지는 선택 사항 — 없으면 http.ErrMissingFile이 반환되며 정상 케이스로 처리한다.
	if fileHeader, err := c.FormFile("profile_image"); err == nil {
		f, openErr := fileHeader.Open()
		if openErr != nil {
			pkg.Fail(c, 400, "이미지 파일을 읽을 수 없습니다")
			return
		}
		defer f.Close()
		req.ProfileImage = f
		req.ProfileImageFilename = fileHeader.Filename
	}

	user, err := h.authSvc.CreateUser(&req)
	if err != nil {
		pkg.Fail(c, 500, err.Error())
		return
	}

	pkg.Success(c, 201, user)
}

func (h *authHandler) CreateOAuthUser(c *gin.Context) {
	req := dto.CreateOAuthUserRequest{
		Username:     c.PostForm("username"),
		Nickname:     c.PostForm("nickname"),
		AgreeTerms:   c.PostForm("agreeTerms") == "true",
		AgreePrivacy: c.PostForm("agreePrivacy") == "true",
	}

	if len(req.Username) < 4 || len(req.Username) > 20 {
		pkg.Fail(c, 400, "username은 4~20자여야 합니다")
		return
	}
	if len(req.Nickname) < 2 || len(req.Nickname) > 20 {
		pkg.Fail(c, 400, "nickname은 2~20자여야 합니다")
		return
	}

	tempToken := c.GetHeader("Authorization")
	if tempToken == "" {
		pkg.Fail(c, 401, "temp token이 없습니다")
		return
	}
	req.TempToken = strings.TrimPrefix(tempToken, "Bearer ")

	// 이미지는 선택 사항 — 없으면 http.ErrMissingFile이 반환되며 정상 케이스로 처리한다.
	if fileHeader, err := c.FormFile("profile_image"); err == nil {
		f, openErr := fileHeader.Open()
		if openErr != nil {
			pkg.Fail(c, 400, "이미지 파일을 읽을 수 없습니다")
			return
		}
		defer f.Close()
		req.ProfileImage = f
		req.ProfileImageFilename = fileHeader.Filename
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

func (h *authHandler) Refresh(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		pkg.Fail(c, 401, "refresh token이 없습니다")
		return
	}
	refreshToken := strings.TrimPrefix(authHeader, "Bearer ")

	res, err := h.authSvc.Refresh(&dto.RefreshRequest{RefreshToken: refreshToken})
	if err != nil {
		pkg.Fail(c, 401, err.Error())
		return
	}

	pkg.Success(c, 200, res)
}
