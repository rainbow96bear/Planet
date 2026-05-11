package handler

import (
	"planet/internal/service"
)

type FeedHandler interface {
}

type feedHandler struct {
	feedSvc service.FeedService
}

func NewFeedHandler(feedSvc service.FeedService) FeedHandler {
	return &feedHandler{feedSvc: feedSvc}
}
