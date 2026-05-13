package service

import (
	"planet/internal/dto"
	"planet/internal/repository"

	"gorm.io/gorm"
)

type FeedService interface {
	GetFeed(userID string) ([]*dto.GetFeedResponse, error)
	GetExploreFeed() ([]*dto.GetFeedResponse, error)
}

type feedService struct {
	db           *gorm.DB
	activityRepo repository.ActivityRepository
}

func NewFeedService(
	db *gorm.DB,
	activityRepo repository.ActivityRepository,
) FeedService {
	return &feedService{
		db:           db,
		activityRepo: activityRepo,
	}
}

func (s *feedService) GetFeed(userID string) ([]*dto.GetFeedResponse, error) {
	return s.activityRepo.FindFeed(userID, 20)
}

func (s *feedService) GetExploreFeed() ([]*dto.GetFeedResponse, error) {
	return s.activityRepo.FindExploreFeed(20)
}
