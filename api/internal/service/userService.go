package service

import (
	"errors"
	"planet/internal/dto"
	"planet/internal/model"
	"planet/internal/repository"

	"gorm.io/gorm"
)

type UserService interface {
	GetProfile(*dto.GetProfileRequest) (*dto.GetProfileResponse, error)
	Follow(*dto.FollowRequest) (*dto.FollowResponse, error)
	Unfollow(*dto.UnfollowRequest) (*dto.UnfollowResponse, error)
}

type userService struct {
	db         *gorm.DB
	userRepo   repository.UserRepository
	followRepo repository.FollowRepository
}

func NewUserService(
	db *gorm.DB,
	userRepo repository.UserRepository,
	followRepo repository.FollowRepository,
) UserService {
	return &userService{
		db:         db,
		userRepo:   userRepo,
		followRepo: followRepo,
	}
}

// service
func (s *userService) GetProfile(req *dto.GetProfileRequest) (*dto.GetProfileResponse, error) {
	user, err := s.userRepo.FindByUserId(req.UserId)
	if err != nil {
		return nil, err
	}

	isFollowing, err := s.followRepo.IsFollowing(req.RequesterUserId, req.UserId)
	if err != nil {
		return nil, err
	}

	return &dto.GetProfileResponse{
		Username:    user.Username,
		Nickname:    user.Nickname,
		IsOwner:     req.UserId == req.RequesterUserId,
		IsFollowing: isFollowing,
	}, nil
}

func (s *userService) Follow(req *dto.FollowRequest) (*dto.FollowResponse, error) {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := s.followRepo.Follow(tx, &model.Follow{
		FollowerID:  req.FollowerID,
		FollowingID: req.FollowingID,
	}); err != nil {
		tx.Rollback()
		return nil, errors.New("이미 팔로우 중입니다")
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &dto.FollowResponse{IsFollowing: true}, nil
}

func (s *userService) Unfollow(req *dto.UnfollowRequest) (*dto.UnfollowResponse, error) {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := s.followRepo.Unfollow(tx, req.FollowerID, req.FollowingID); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &dto.UnfollowResponse{IsFollowing: false}, nil
}
