package handler

import (
	"planet/internal/pkg"
	"planet/internal/service"

	"github.com/gin-gonic/gin"
)

type FeedHandler interface {
	GetFeed(c *gin.Context)
	GetExploreFeed(c *gin.Context)
}

type feedHandler struct {
	feedSvc service.FeedService
}

func NewFeedHandler(feedSvc service.FeedService) FeedHandler {
	return &feedHandler{feedSvc: feedSvc}
}

func (h *feedHandler) GetFeed(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		// 비로그인이면 explore로
		h.GetExploreFeed(c)
		return
	}

	feed, err := h.feedSvc.GetFeed(userID)
	if err != nil {
		pkg.Fail(c, 500, "피드를 불러오지 못했습니다")
		return
	}

	pkg.Success(c, 200, feed)
}

func (h *feedHandler) GetExploreFeed(c *gin.Context) {
	feed, err := h.feedSvc.GetExploreFeed()
	if err != nil {
		pkg.Fail(c, 500, "피드를 불러오지 못했습니다")
		return
	}

	pkg.Success(c, 200, feed)
}
