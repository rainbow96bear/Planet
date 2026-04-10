package service

import (
	"planet/internal/dto"
	"planet/internal/model"
	"planet/internal/repository"
)

type UserService interface {
	CreateUser(*dto.CreateUserRequest) error
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) CreateUser(req *dto.CreateUserRequest) error {
	user := &model.User{
		Username: req.Username,
		Nickname: req.Nickname,
		Password: req.Password,
	}
	return s.userRepo.CreateWithTx(user)
}
