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
	db        *gorm.DB
	userRepo  repository.UserRepository
	orbitRepo repository.OrbitRepository
}

func NewSearchService(
	db *gorm.DB,
	userRepo repository.UserRepository,
	orbitRepo repository.OrbitRepository,
) SearchService {
	return &searchService{
		db:        db,
		userRepo:  userRepo,
		orbitRepo: orbitRepo,
	}
}

func (s *searchService) SearchUsers(req *dto.SearchUsersRequest) ([]*dto.SearchUsersResponse, error) {
	users, err := s.userRepo.SearchByKeyword(req.Q)
	if err != nil {
		return nil, err
	}

	result := make([]*dto.SearchUsersResponse, len(users))
	for i, u := range users {
		isOrbiting := false
		if req.RequesterUserId != "" {
			isOrbiting, err = s.orbitRepo.IsOrbiting(req.RequesterUserId, u.ID)
			if err != nil {
				return nil, err
			}
		}

		result[i] = &dto.SearchUsersResponse{
			UserId:     u.ID,
			Username:   u.Username,
			Nickname:   u.Nickname,
			IsOrbiting: isOrbiting,
		}
	}
	return result, nil
}
