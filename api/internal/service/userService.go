package service

import (
	"planet/internal/dto"
	"planet/internal/model"
	"planet/internal/repository"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService interface {
	CreateUser(*dto.CreateUserRequest) (*dto.CreateUserResponse, error)
}

type userService struct {
	db       *gorm.DB
	userRepo repository.UserRepository
}

func NewUserService(db *gorm.DB, userRepo repository.UserRepository) UserService {
	return &userService{
		db:       db,
		userRepo: userRepo,
	}
}

func (s *userService) CreateUser(req *dto.CreateUserRequest) (*dto.CreateUserResponse, error) {
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
	if err := s.userRepo.Create(tx, user); err != nil {
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
