package service

import (
	"planet/internal/dto"
	"planet/internal/repository"

	"gorm.io/gorm"
)

type FeedService interface {
	GetFeed(userID string) ([]*dto.GetFeedResponse, error)
	GetExploreFeed(userID string) ([]*dto.GetFeedResponse, error)
}

type feedService struct {
	db       *gorm.DB
	feedRepo repository.FeedRepository
}

func NewFeedService(
	db *gorm.DB,
	feedRepo repository.FeedRepository,
) FeedService {
	return &feedService{
		db:       db,
		feedRepo: feedRepo,
	}
}

func (s *feedService) GetFeed(userID string) ([]*dto.GetFeedResponse, error) {
	return s.feedRepo.FindFeed(userID, 20)
}

func (s *feedService) GetExploreFeed(userID string) ([]*dto.GetFeedResponse, error) {
	return s.feedRepo.FindExploreFeed(userID, 20)
}
