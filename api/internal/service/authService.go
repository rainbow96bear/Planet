package service

import (
	"errors"
	"io"
	"log"
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
		return nil, err
	}

	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

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

	if err := s.userRepo.CreateUser(tx, user); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

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

	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	claims, err := pkg.ParseTempToken(req.TempToken)
	if err != nil {
		tx.Rollback()
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

	if err := s.userRepo.CreateUser(tx, user); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

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
		log.Printf("profile image upload failed for user %v: %v", user.ID, err)
		return
	}

	if err := s.userRepo.UpdateProfileImage(s.db, user.ID, url); err != nil {
		log.Printf("failed to save profile image url for user %v: %v", user.ID, err)
		return
	}

	user.ProfileImage = url
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

	if user.UserStatus != model.UserStatusActive {
		return nil, errors.New("account unavailable")
	}

	// 저장된 hash와 입력된 password 비교
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid password")
	}

	now := time.Now()
	if err := s.userRepo.UpdateLastLogin(user.ID, now); err != nil {
		// 로그인을 실패시키진 않고 로그만 남기는 걸 추천
		log.Printf("failed to update last login: %v", err)
	}

	accessToken, err := pkg.GenerateAccessToken(user.ID, user.Username)
	if err != nil {

	}
	refreshToken, err := pkg.GenerateRefreshToken(user.ID, user.Username)
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
	log.Printf("provider : %+v, provider id: %+v\n", req.Provider, req.ProviderID)
	if err != nil {
		log.Printf("FindByProviderInfo err : %+v\n", err.Error())
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
	accessToken, err := pkg.GenerateAccessToken(user.ID, user.Username)
	if err != nil {
		return nil, err
	}
	refreshToken, err := pkg.GenerateRefreshToken(user.ID, user.Username)
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
	_, err = s.userRepo.FindByUsername(username)
	if err != nil {
		return nil, errors.New("user not found")
	}

	accessToken, err := pkg.GenerateAccessToken(userid, username)
	if err != nil {
		return nil, err
	}
	refreshToken, err := pkg.GenerateRefreshToken(userid, username)
	if err != nil {
		return nil, err
	}

	return &dto.RefreshResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
