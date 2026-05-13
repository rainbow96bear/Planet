package handler

import (
	"planet/internal/pkg"
	"planet/internal/service"

	"github.com/gin-gonic/gin"
)

type NotificationHandler interface {
	GetNotifications(c *gin.Context)
	GetUnreadCount(c *gin.Context)
	MarkAllAsRead(c *gin.Context)
}

type notificationHandler struct {
	notificationSvc service.NotificationService
}

func NewNotificationHandler(notificationSvc service.NotificationService) NotificationHandler {
	return &notificationHandler{notificationSvc: notificationSvc}
}

func (h *notificationHandler) GetNotifications(c *gin.Context) {
	userID := c.GetString("userID")

	notifications, err := h.notificationSvc.GetNotifications(userID)
	if err != nil {
		pkg.Fail(c, 500, "알림을 불러오지 못했습니다")
		return
	}

	pkg.Success(c, 200, notifications)
}

func (h *notificationHandler) GetUnreadCount(c *gin.Context) {
	userID := c.GetString("userID")

	result, err := h.notificationSvc.GetUnreadCount(userID)
	if err != nil {
		pkg.Fail(c, 500, "알림 수를 불러오지 못했습니다")
		return
	}

	pkg.Success(c, 200, result)
}

func (h *notificationHandler) MarkAllAsRead(c *gin.Context) {
	userID := c.GetString("userID")

	if err := h.notificationSvc.MarkAllAsRead(userID); err != nil {
		pkg.Fail(c, 500, "읽음 처리에 실패했습니다")
		return
	}

	pkg.Success(c, 200, nil)
}
