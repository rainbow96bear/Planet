package service

import (
	"planet/internal/repository"

	"gorm.io/gorm"
)

type FeedService interface {
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
