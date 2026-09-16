package service

import (
	"errors"
	"io"
	"log/slog"
	"planet/internal/constants"
	"planet/internal/dto"
	"planet/internal/model"
	"planet/internal/pkg"
	"planet/internal/repository"
	"planet/internal/storage"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService interface {
	CreateUser(*dto.CreateUserRequest) (*dto.CreateUserResponse, error)
	CreateOAuthUser(*dto.CreateOAuthUserRequest) (*dto.CreateOAuthUserResponse, error)
	IsUsernameAvailable(*dto.CheckUsernameRequest) (*dto.CheckUsernameResponse, error)
	Login(*dto.LoginRequest) (*dto.LoginResponse, error)
	OauthLogin(*dto.OauthLoginRequest) (*dto.OauthLoginResponse, error)
	Refresh(*dto.RefreshRequest) (*dto.RefreshResponse, error)
}

type authService struct {
	db          *gorm.DB
	userRepo    repository.UserRepository
	fileStorage storage.FileStorage
}

func NewAuthService(db *gorm.DB, userRepo repository.UserRepository, fileStorage storage.FileStorage) AuthService {
	return &authService{
		db:          db,
		userRepo:    userRepo,
		fileStorage: fileStorage,
	}
}

func (s *authService) CreateUser(req *dto.CreateUserRequest) (*dto.CreateUserResponse, error) {
	if !req.AgreeTerms {
		return nil, errors.New("terms agreement is required")
	}

	if !req.AgreePrivacy {
		return nil, errors.New("privacy agreement is required")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("password hashing failed", "username", req.Username, "error", err)
		return nil, err
	}

	now := time.Now()

	user := &model.User{
		Username: req.Username,
		Nickname: req.Nickname,
		Password: string(hashed),

		Provider: "local",

		TermsVersion:   constants.CurrentTermsVersion,
		PrivacyVersion: constants.CurrentPrivacyVersion,

		TermsAgreedAt:   &now,
		PrivacyAgreedAt: &now,
	}

	if err := s.db.Transaction(func(tx *gorm.DB) error {
		return s.userRepo.CreateUser(tx, user)
	}); err != nil {
		slog.Error("failed to create user", "username", req.Username, "error", err)
		return nil, err
	}

	slog.Info("user created", "user_id", user.ID, "username", user.Username, "provider", user.Provider)

	// 이미지 업로드는 DB 트랜잭션 밖에서 best-effort로 처리 (실패해도 가입은 유지)
	s.attachProfileImageBestEffort(user, req.ProfileImage, req.ProfileImageFilename)

	return &dto.CreateUserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Nickname:  user.Nickname,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (s *authService) CreateOAuthUser(req *dto.CreateOAuthUserRequest) (*dto.CreateOAuthUserResponse, error) {
	if !req.AgreeTerms {
		return nil, errors.New("terms agreement is required")
	}

	if !req.AgreePrivacy {
		return nil, errors.New("privacy agreement is required")
	}

	// DB와 무관한 순수 검증이라 트랜잭션 밖에서 먼저 처리한다.
	claims, err := pkg.ParseTempToken(req.TempToken)
	if err != nil {
		slog.Warn("invalid temp token on oauth signup", "error", err)
		return nil, errors.New("invalid temp token")
	}

	now := time.Now()

	user := &model.User{
		Username:   req.Username,
		Nickname:   req.Nickname,
		Provider:   claims.Provider,
		ProviderID: claims.ProviderID,

		TermsVersion:   constants.CurrentTermsVersion,
		PrivacyVersion: constants.CurrentPrivacyVersion,

		TermsAgreedAt:   &now,
		PrivacyAgreedAt: &now,
	}

	if err := s.db.Transaction(func(tx *gorm.DB) error {
		return s.userRepo.CreateUser(tx, user)
	}); err != nil {
		slog.Error("failed to create oauth user", "username", req.Username, "provider", claims.Provider, "error", err)
		return nil, err
	}

	slog.Info("oauth user created", "user_id", user.ID, "username", user.Username, "provider", user.Provider)

	// 이미지 업로드는 DB 트랜잭션 밖에서 best-effort로 처리 (실패해도 가입은 유지)
	s.attachProfileImageBestEffort(user, req.ProfileImage, req.ProfileImageFilename)

	return &dto.CreateOAuthUserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Nickname:  user.Nickname,
		CreatedAt: user.CreatedAt,
	}, nil
}

// attachProfileImageBestEffort는 이미지가 있으면 업로드해서 유저의 ProfileImage를 갱신한다.
// 이미지 업로드는 회원가입의 필수 조건이 아니므로, 실패해도 에러를 반환하지 않고
// 로그만 남긴다 (가입 자체를 막지 않기 위한 best-effort 정책).
func (s *authService) attachProfileImageBestEffort(user *model.User, image io.Reader, filename string) {
	if image == nil {
		return
	}

	key := storage.ProfileImageKey(user.ID, filename)
	url, err := s.fileStorage.Upload(storage.UploadInput{
		Key:    key,
		Reader: image,
	})
	if err != nil {
		slog.Warn("profile image upload failed", "user_id", user.ID, "error", err)
		return
	}

	if err := s.userRepo.UpdateProfileImage(s.db, user.ID, url); err != nil {
		slog.Warn("failed to save profile image url", "user_id", user.ID, "error", err)
		return
	}

	user.ProfileImage = url
	slog.Info("profile image attached", "user_id", user.ID)
}

func (s *authService) IsUsernameAvailable(req *dto.CheckUsernameRequest) (*dto.CheckUsernameResponse, error) {
	exists, err := s.userRepo.IsUsernameExists(req.Username)
	if err != nil {
		slog.Error("username availability check failed", "username", req.Username, "error", err)
		return nil, err
	}

	return &dto.CheckUsernameResponse{
		Username:  req.Username,
		Available: !exists,
	}, nil
}

func (s *authService) Login(req *dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		slog.Warn("login failed: user not found", "username", req.Username)
		return nil, errors.New("user not found")
	}

	if user.UserStatus != model.UserStatusActive {
		slog.Warn("login failed: inactive account", "user_id", user.ID, "status", user.UserStatus)
		return nil, errors.New("account unavailable")
	}

	// 저장된 hash와 입력된 password 비교
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		slog.Warn("login failed: invalid password", "user_id", user.ID)
		return nil, errors.New("invalid password")
	}

	now := time.Now()
	if err := s.userRepo.UpdateLastLogin(user.ID, now); err != nil {
		// 로그인을 실패시키진 않고 로그만 남긴다
		slog.Warn("failed to update last login", "user_id", user.ID, "error", err)
	}

	accessToken, err := pkg.GenerateAccessToken(user.ID, user.Username)
	if err != nil {
		slog.Error("access token generation failed", "user_id", user.ID, "error", err)
		return nil, err
	}
	refreshToken, err := pkg.GenerateRefreshToken(user.ID, user.Username)
	if err != nil {
		slog.Error("refresh token generation failed", "user_id", user.ID, "error", err)
		return nil, err
	}

	slog.Info("login success", "user_id", user.ID, "username", user.Username)

	return &dto.LoginResponse{
		Username:     user.Username,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *authService) OauthLogin(req *dto.OauthLoginRequest) (*dto.OauthLoginResponse, error) {
	user, err := s.userRepo.FindByProviderInfo(req.Provider, req.ProviderID)
	if err != nil {
		slog.Info("oauth login: new user, issuing temp token", "provider", req.Provider, "provider_id", req.ProviderID)

		tempToken, err := pkg.GenerateTempToken(req.Provider, req.ProviderID)
		if err != nil {
			slog.Error("temp token generation failed", "provider", req.Provider, "error", err)
			return nil, err
		}
		return &dto.OauthLoginResponse{
			IsNewUser: true,
			TempToken: tempToken,
		}, nil
	}

	// 기존 유저면 JWT 발급
	accessToken, err := pkg.GenerateAccessToken(user.ID, user.Username)
	if err != nil {
		slog.Error("access token generation failed", "user_id", user.ID, "error", err)
		return nil, err
	}
	refreshToken, err := pkg.GenerateRefreshToken(user.ID, user.Username)
	if err != nil {
		slog.Error("refresh token generation failed", "user_id", user.ID, "error", err)
		return nil, err
	}

	slog.Info("oauth login success", "user_id", user.ID, "username", user.Username, "provider", req.Provider)

	return &dto.OauthLoginResponse{
		IsNewUser:    false,
		Username:     user.Username,
		Nickname:     user.Nickname,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *authService) Refresh(req *dto.RefreshRequest) (*dto.RefreshResponse, error) {
	claims, err := pkg.ParseRefreshToken(req.RefreshToken)
	if err != nil {
		slog.Warn("refresh failed: invalid refresh token", "error", err)
		return nil, err
	}

	userid := claims.UserID
	username := claims.Username
	if _, err := s.userRepo.FindByUsername(username); err != nil {
		slog.Warn("refresh failed: user not found", "username", username)
		return nil, errors.New("user not found")
	}

	accessToken, err := pkg.GenerateAccessToken(userid, username)
	if err != nil {
		slog.Error("access token generation failed", "user_id", userid, "error", err)
		return nil, err
	}
	refreshToken, err := pkg.GenerateRefreshToken(userid, username)
	if err != nil {
		slog.Error("refresh token generation failed", "user_id", userid, "error", err)
		return nil, err
	}

	slog.Info("token refreshed", "user_id", userid, "username", username)

	return &dto.RefreshResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
