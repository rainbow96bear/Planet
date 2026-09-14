package service

import (
	"log/slog"
	"planet/internal/dto"
	"planet/internal/repository"

	"gorm.io/gorm"
)

type NotificationService interface {
	GetNotifications(receiverID string) ([]*dto.GetNotificationsResponse, error)
	GetUnreadCount(receiverID string) (*dto.GetUnreadCountResponse, error)
	MarkAllAsRead(receiverID string) error
}

type notificationService struct {
	db               *gorm.DB
	notificationRepo repository.NotificationRepository
	userRepo         repository.UserRepository
}

func NewNotificationService(
	db *gorm.DB,
	notificationRepo repository.NotificationRepository,
	userRepo repository.UserRepository,
) NotificationService {
	return &notificationService{
		db:               db,
		notificationRepo: notificationRepo,
		userRepo:         userRepo,
	}
}

func (s *notificationService) GetNotifications(receiverID string) ([]*dto.GetNotificationsResponse, error) {
	notifications, err := s.notificationRepo.FindByReceiverID(receiverID)
	if err != nil {
		slog.Error("failed to fetch notifications", "receiver_id", receiverID, "error", err)
		return nil, err
	}

	result := make([]*dto.GetNotificationsResponse, 0, len(notifications))
	for _, n := range notifications {
		actor, err := s.userRepo.FindByUserId(n.ActorID)
		if err != nil {
			// actor 조회 실패는 알림 하나만 건너뛰고 전체 목록은 계속 반환한다.
			// 다만 원인 파악을 위해(탈퇴 유저인지, DB 문제인지) 흔적은 남긴다.
			slog.Warn("dropped notification: actor not found",
				"notification_id", n.ID,
				"actor_id", n.ActorID,
				"error", err,
			)
			continue
		}
		result = append(result, &dto.GetNotificationsResponse{
			ID:            n.ID,
			ActorID:       n.ActorID,
			ActorNickname: actor.Nickname,
			Type:          string(n.Type),
			IsRead:        n.IsRead,
			CreatedAt:     n.CreatedAt,
		})
	}
	return result, nil
}

func (s *notificationService) GetUnreadCount(receiverID string) (*dto.GetUnreadCountResponse, error) {
	count, err := s.notificationRepo.CountUnread(receiverID)
	if err != nil {
		slog.Error("failed to count unread notifications", "receiver_id", receiverID, "error", err)
		return nil, err
	}
	return &dto.GetUnreadCountResponse{Count: count}, nil
}

func (s *notificationService) MarkAllAsRead(receiverID string) error {
	if err := s.notificationRepo.MarkAllAsRead(receiverID); err != nil {
		slog.Error("failed to mark notifications as read", "receiver_id", receiverID, "error", err)
		return err
	}
	return nil
}
