package service

import (
	"errors"
	"log/slog"
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
	taskRepo         repository.TaskRepository
	notificationRepo repository.NotificationRepository
}

func NewReactionService(
	db *gorm.DB,
	reactionRepo repository.ReactionRepository,
	taskRepo repository.TaskRepository,
	notificationRepo repository.NotificationRepository,
) ReactionService {
	return &reactionService{
		db:               db,
		reactionRepo:     reactionRepo,
		taskRepo:         taskRepo,
		notificationRepo: notificationRepo,
	}
}

func (s *reactionService) AddReaction(req *dto.AddReactionRequest) (*dto.AddReactionResponse, error) {
	reactionType := model.ReactionType(req.Type)
	if reactionType != model.ReactionTypeLike && reactionType != model.ReactionTypeCheer {
		return nil, errors.New("invalid reaction type")
	}

	// task 조회해서 owner 찾기
	task, err := s.taskRepo.GetTaskByID(req.TaskID)
	if err != nil {
		return nil, errors.New("존재하지 않는 할 일입니다")
	}

	var response *dto.AddReactionResponse
	err = s.db.Transaction(func(tx *gorm.DB) error {
		reaction := &model.Reaction{
			TaskID: req.TaskID,
			UserID: req.UserID,
			Type:   reactionType,
		}
		if err := s.reactionRepo.Upsert(tx, reaction); err != nil {
			slog.Error("failed to upsert reaction",
				"task_id", req.TaskID,
				"user_id", req.UserID,
				"type", req.Type,
				"error", err,
			)
			return err
		}

		// TaskOwnerID 대신 task.UserID 직접 사용
		if task.UserID != req.UserID {
			n := &model.Notification{
				ReceiverID: task.UserID,
				ActorID:    req.UserID,
				TargetID:   req.TaskID,
				TargetType: model.NotificationTargetTypeTask,
				Type:       model.NotificationTypeReaction,
				IsRead:     false,
			}
			if err := s.notificationRepo.Upsert(tx, n); err != nil {
				slog.Error("failed to create reaction notification",
					"task_id", req.TaskID,
					"receiver_id", task.UserID,
					"actor_id", req.UserID,
					"error", err,
				)
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
		if err := s.reactionRepo.Delete(tx, req.TaskID, req.UserID, reactionType); err != nil {
			slog.Error("failed to delete reaction",
				"task_id", req.TaskID,
				"user_id", req.UserID,
				"type", req.Type,
				"error", err,
			)
			return err
		}
		if err := s.notificationRepo.DeleteByActorAndTask(tx, req.UserID, req.TaskID, model.NotificationTypeReaction); err != nil {
			slog.Error("failed to delete reaction notification",
				"task_id", req.TaskID,
				"actor_id", req.UserID,
				"error", err,
			)
			return err
		}
		return nil
	})
}
