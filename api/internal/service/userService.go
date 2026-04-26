package service

import (
	"planet/internal/dto"
	"planet/internal/repository"

	"gorm.io/gorm"
)

type UserService interface {
	GetProfile(*dto.GetProfileRequest) (*dto.GetProfileResponse, error)
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

// service
func (s *userService) GetProfile(req *dto.GetProfileRequest) (*dto.GetProfileResponse, error) {
	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		return nil, err
	}

	return &dto.GetProfileResponse{
		Username: user.Username,
		Nickname: user.Nickname,
		IsOwner:  req.Username == req.RequesterUsername,
	}, nil
}
