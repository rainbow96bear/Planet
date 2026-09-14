package service

import (
	"log/slog"
	"planet/internal/dto"
	"planet/internal/repository"

	"gorm.io/gorm"
)

const feedPageSize = 20

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
	feed, err := s.feedRepo.FindFeed(userID, feedPageSize)
	if err != nil {
		slog.Error("failed to fetch feed", "user_id", userID, "error", err)
		return nil, err
	}
	return feed, nil
}

func (s *feedService) GetExploreFeed(userID string) ([]*dto.GetFeedResponse, error) {
	feed, err := s.feedRepo.FindExploreFeed(userID, feedPageSize)
	if err != nil {
		slog.Error("failed to fetch explore feed", "user_id", userID, "error", err)
		return nil, err
	}
	return feed, nil
}
