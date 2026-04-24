package service

import (
	"errors"
	"planet/internal/dto"
	"planet/internal/model"
	"planet/internal/pkg"
	"planet/internal/repository"

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
	db       *gorm.DB
	userRepo repository.UserRepository
}

func NewAuthService(db *gorm.DB, userRepo repository.UserRepository) AuthService {
	return &authService{
		db:       db,
		userRepo: userRepo,
	}
}

func (s *authService) CreateUser(req *dto.CreateUserRequest) (*dto.CreateUserResponse, error) {
	// 1. 비밀번호 해싱
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 2. tx 시작
	tx := s.db.Begin()

	// panic 대비
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	user := &model.User{
		Username: req.Username,
		Nickname: req.Nickname,
		Password: string(hashed),
	}

	// 3. 저장
	if err := s.userRepo.CreateUser(tx, user); err != nil {
		tx.Rollback()

		return nil, err
	}

	// 4. commit
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// 5. DTO 반환
	return &dto.CreateUserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Nickname:  user.Nickname,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (s *authService) CreateOAuthUser(req *dto.CreateOAuthUserRequest) (*dto.CreateOAuthUserResponse, error) {
	// tx 시작
	tx := s.db.Begin()

	// panic 대비
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	claims, err := pkg.ParseTempToken(req.TempToken)
	if err != nil {
		return nil, errors.New("invalid temp token")
	}

	user := &model.User{
		Username:   req.Username,
		Nickname:   req.Nickname,
		Provider:   claims.Provider,
		ProviderID: claims.ProviderID,
	}

	// 저장
	if err := s.userRepo.CreateUser(tx, user); err != nil {
		tx.Rollback()

		return nil, err
	}

	// commit
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// DTO 반환
	return &dto.CreateOAuthUserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Nickname:  user.Nickname,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (s *authService) IsUsernameAvailable(req *dto.CheckUsernameRequest) (*dto.CheckUsernameResponse, error) {
	exists, err := s.userRepo.IsUsernameExists(req.Username)
	if err != nil {
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
		return nil, errors.New("user not found")
	}

	// 저장된 hash와 입력된 password 비교
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid password")
	}

	accessToken, err := pkg.GenerateAccessToken(user.ID, user.Username, user.Nickname)
	if err != nil {

	}
	refreshToken, err := pkg.GenerateRefreshToken(user.ID, user.Username, user.Nickname)
	if err != nil {

	}
	return &dto.LoginResponse{
		Username:     user.Username,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *authService) OauthLogin(req *dto.OauthLoginRequest) (*dto.OauthLoginResponse, error) {
	user, err := s.userRepo.FindByProviderInfo(req.Provider, req.ProviderID)
	if err != nil {
		// 유저 없으면 temp_token 발급
		tempToken, err := pkg.GenerateTempToken(req.Provider, req.ProviderID)
		if err != nil {
			return nil, err
		}
		return &dto.OauthLoginResponse{
			IsNewUser: true,
			TempToken: tempToken,
		}, nil
	}

	// 기존 유저면 JWT 발급
	accessToken, err := pkg.GenerateAccessToken(user.ID, user.Username, user.Nickname)
	if err != nil {
		return nil, err
	}
	refreshToken, err := pkg.GenerateRefreshToken(user.ID, user.Username, user.Nickname)
	if err != nil {
		return nil, err
	}

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
		return nil, err
	}

	userid := claims.UserID
	username := claims.Username
	nickname := claims.Nickname
	_, err = s.userRepo.FindByUsername(username)
	if err != nil {
		return nil, errors.New("user not found")
	}

	accessToken, err := pkg.GenerateAccessToken(userid, username, nickname)
	if err != nil {
		return nil, err
	}
	refreshToken, err := pkg.GenerateRefreshToken(userid, username, nickname)
	if err != nil {
		return nil, err
	}

	return &dto.RefreshResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
