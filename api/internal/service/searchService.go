package service

import (
	"planet/internal/dto"
	"planet/internal/repository"

	"gorm.io/gorm"
)

type SearchService interface {
	SearchUsers(*dto.SearchUsersRequest) ([]*dto.SearchUsersResponse, error)
}

type searchService struct {
	db         *gorm.DB
	userRepo   repository.UserRepository
	followRepo repository.FollowRepository
}

func NewSearchService(
	db *gorm.DB,
	userRepo repository.UserRepository,
	followRepo repository.FollowRepository,
) SearchService {
	return &searchService{
		db:         db,
		userRepo:   userRepo,
		followRepo: followRepo,
	}
}

func (s *searchService) SearchUsers(req *dto.SearchUsersRequest) ([]*dto.SearchUsersResponse, error) {
	users, err := s.userRepo.SearchByKeyword(req.Q)
	if err != nil {
		return nil, err
	}

	result := make([]*dto.SearchUsersResponse, len(users))
	for i, u := range users {
		isFollowing := false
		if req.RequesterUserId != 0 {
			isFollowing, err = s.followRepo.IsFollowing(req.RequesterUserId, u.ID)
			if err != nil {
				return nil, err
			}
		}

		result[i] = &dto.SearchUsersResponse{
			UserId:      u.ID,
			Username:    u.Username,
			Nickname:    u.Nickname,
			IsFollowing: isFollowing,
		}
	}
	return result, nil
}
