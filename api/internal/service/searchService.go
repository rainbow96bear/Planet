package service

import (
	"log/slog"
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
		slog.Error("failed to search users", "query", req.Q, "error", err)
		return nil, err
	}

	// NOTE: 결과 수만큼 IsOrbiting을 개별 호출해서 N+1 쿼리가 발생한다.
	// 트래픽/결과 수가 늘어나면 orbitRepo에 배치 조회(IsOrbitingBatch 등)를 추가하는 걸 고려할 것.
	result := make([]*dto.SearchUsersResponse, len(users))
	for i, u := range users {
		isOrbiting := false
		if req.RequesterUserId != "" {
			isOrbiting, err = s.orbitRepo.IsOrbiting(req.RequesterUserId, u.ID)
			if err != nil {
				slog.Error("failed to check orbit status",
					"requester_id", req.RequesterUserId,
					"target_id", u.ID,
					"error", err,
				)
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
