package service

import (
	"planet/internal/dto"
	"planet/internal/repository"
)

type NotificationService interface {
	GetNotifications(receiverID string) ([]*dto.GetNotificationsResponse, error)
	GetUnreadCount(receiverID string) (*dto.GetUnreadCountResponse, error)
	MarkAllAsRead(receiverID string) error
}

type notificationService struct {
	notificationRepo repository.NotificationRepository
	userRepo         repository.UserRepository
}

func NewNotificationService(
	notificationRepo repository.NotificationRepository,
	userRepo repository.UserRepository,
) NotificationService {
	return &notificationService{
		notificationRepo: notificationRepo,
		userRepo:         userRepo,
	}
}

func (s *notificationService) GetNotifications(receiverID string) ([]*dto.GetNotificationsResponse, error) {
	notifications, err := s.notificationRepo.FindByReceiverID(receiverID)
	if err != nil {
		return nil, err
	}

	result := make([]*dto.GetNotificationsResponse, 0, len(notifications))
	for _, n := range notifications {
		actor, err := s.userRepo.FindByUserId(n.ActorID)
		if err != nil {
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
		return nil, err
	}
	return &dto.GetUnreadCountResponse{Count: count}, nil
}

func (s *notificationService) MarkAllAsRead(receiverID string) error {
	return s.notificationRepo.MarkAllAsRead(receiverID)
}
