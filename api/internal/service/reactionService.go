package service

import (
	"errors"
	"planet/internal/dto"
	"planet/internal/model"
	"planet/internal/repository"

	"gorm.io/gorm"
)

type ReactionService interface {
	AddReaction(req *dto.AddReactionRequest) (*dto.AddReactionResponse, error)
	RemoveReaction(req *dto.RemoveReactionRequest) error
}

type reactionService struct {
	db               *gorm.DB
	reactionRepo     repository.ReactionRepository
	notificationRepo repository.NotificationRepository
}

func NewReactionService(
	db *gorm.DB,
	reactionRepo repository.ReactionRepository,
	notificationRepo repository.NotificationRepository,
) ReactionService {
	return &reactionService{
		db:               db,
		reactionRepo:     reactionRepo,
		notificationRepo: notificationRepo,
	}
}

func (s *reactionService) AddReaction(req *dto.AddReactionRequest) (*dto.AddReactionResponse, error) {
	var response *dto.AddReactionResponse

	err := s.db.Transaction(func(tx *gorm.DB) error {
		reaction := &model.Reaction{
			TaskID: req.TaskID,
			UserID: req.UserID,
			Type:   req.Type,
		}
		if err := s.reactionRepo.Upsert(reaction); err != nil {
			return err
		}

		// 자기 자신 태스크엔 알림 안 보냄
		if req.TaskOwnerID != req.UserID {
			n := &model.Notification{
				ReceiverID: req.TaskOwnerID,
				ActorID:    req.UserID,
				TargetID:   req.TaskID,
				TargetType: model.NotificationTargetTypeTask,
				Type:       model.NotificationTypeReaction,
				IsRead:     false,
			}
			if err := s.notificationRepo.Upsert(tx, n); err != nil {
				return err
			}
		}

		response = &dto.AddReactionResponse{
			TaskID: req.TaskID,
			Type:   req.Type,
		}
		return nil
	})

	return response, err
}

func (s *reactionService) RemoveReaction(req *dto.RemoveReactionRequest) error {
	reactionType := model.ReactionType(req.Type)
	if reactionType != model.ReactionTypeLike && reactionType != model.ReactionTypeCheer {
		return errors.New("invalid reaction type")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.reactionRepo.Delete(req.TaskID, req.UserID, reactionType); err != nil {
			return err
		}

		return s.notificationRepo.DeleteByActorAndTask(
			tx,
			req.UserID,
			req.TaskID,
			model.NotificationTypeReaction,
		)
	})
}
