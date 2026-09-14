package service

import (
	"errors"
	"log/slog"
	"planet/internal/dto"
	"planet/internal/model"
	"planet/internal/repository"

	"gorm.io/gorm"
)

type TaskService interface {
	CreateTask(*dto.CreateTaskRequest) (*dto.CreateTaskResponse, error)
	DeleteTask(*dto.DeleteTaskRequest) error
	GetTasksByMonth(*dto.GetTasksByMonthRequest) ([]*dto.GetTasksByMonthResponse, error)
	ToggleTask(*dto.ToggleTaskRequest) (*dto.ToggleTaskResponse, error)
}

type taskService struct {
	db           *gorm.DB
	taskRepo     repository.TaskRepository
	feedRepo     repository.FeedRepository
	reactionRepo repository.ReactionRepository
}

func NewTaskService(
	db *gorm.DB,
	taskRepo repository.TaskRepository,
	feedRepo repository.FeedRepository,
	reactionRepo repository.ReactionRepository,
) TaskService {
	return &taskService{
		db:           db,
		taskRepo:     taskRepo,
		feedRepo:     feedRepo,
		reactionRepo: reactionRepo,
	}
}

func (s *taskService) CreateTask(req *dto.CreateTaskRequest) (*dto.CreateTaskResponse, error) {
	tx := s.db.Begin()
	if tx.Error != nil {
		slog.Error("failed to begin transaction", "user_id", req.UserID, "error", tx.Error)
		return nil, tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	task := &model.Task{
		UserID:      req.UserID,
		Title:       req.Title,
		Description: req.Description,
		Date:        req.Date,
		IsPublic:    req.IsPublic,
	}
	if err := s.taskRepo.CreateTask(tx, task); err != nil {
		tx.Rollback()
		slog.Error("failed to create task", "user_id", req.UserID, "error", err)
		return nil, err
	}

	// TargetID/TargetType → TaskID, UserID → ActorID
	if err := s.feedRepo.Create(tx, &model.Feed{
		ActorID: req.UserID,
		TaskID:  task.ID,
		Type:    model.TaskCreated,
	}); err != nil {
		tx.Rollback()
		slog.Error("failed to create feed for task", "task_id", task.ID, "user_id", req.UserID, "error", err)
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		slog.Error("failed to commit task creation", "user_id", req.UserID, "error", err)
		return nil, err
	}

	slog.Info("task created", "task_id", task.ID, "user_id", req.UserID, "is_public", task.IsPublic)

	return &dto.CreateTaskResponse{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Date:        task.Date,
		IsCompleted: task.IsCompleted,
		IsPublic:    task.IsPublic,
	}, nil
}

func (s *taskService) DeleteTask(req *dto.DeleteTaskRequest) error {
	task, err := s.taskRepo.GetTaskByID(req.ID)
	if err != nil {
		return errors.New("존재하지 않는 할 일입니다")
	}
	if task.UserID != req.UserID {
		slog.Warn("unauthorized task delete attempt",
			"task_id", req.ID,
			"owner_id", task.UserID,
			"requester_id", req.UserID,
		)
		return errors.New("권한이 없습니다")
	}

	tx := s.db.Begin()
	if tx.Error != nil {
		slog.Error("failed to begin transaction", "task_id", req.ID, "error", tx.Error)
		return tx.Error
	}

	if err := s.taskRepo.DeleteTask(tx, req.ID); err != nil {
		tx.Rollback()
		slog.Error("failed to delete task", "task_id", req.ID, "user_id", req.UserID, "error", err)
		return err
	}

	if err := s.feedRepo.DeleteByTaskID(tx, req.ID); err != nil {
		tx.Rollback()
		slog.Error("failed to delete feed for task", "task_id", req.ID, "error", err)
		return err
	}

	if err := s.reactionRepo.DeleteByTaskID(tx, req.ID); err != nil {
		tx.Rollback()
		slog.Error("failed to delete reactions for task", "task_id", req.ID, "error", err)
		return err
	}

	if err := tx.Commit().Error; err != nil {
		slog.Error("failed to commit task deletion", "task_id", req.ID, "error", err)
		return err
	}

	slog.Info("task deleted", "task_id", req.ID, "user_id", req.UserID)
	return nil
}

func (s *taskService) GetTasksByMonth(req *dto.GetTasksByMonthRequest) ([]*dto.GetTasksByMonthResponse, error) {
	isOwner := req.UserID == req.RequesterUserId

	tasks, err := s.taskRepo.GetTasksByMonth(req.UserID, req.Year, req.Month, isOwner)
	if err != nil {
		slog.Error("failed to fetch tasks by month",
			"user_id", req.UserID,
			"year", req.Year,
			"month", req.Month,
			"error", err,
		)
		return nil, err
	}

	result := make([]*dto.GetTasksByMonthResponse, len(tasks))
	for i, task := range tasks {
		result[i] = &dto.GetTasksByMonthResponse{
			ID:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			Date:        task.Date,
			IsCompleted: task.IsCompleted,
			IsPublic:    task.IsPublic,
		}
	}
	return result, nil
}

func (s *taskService) ToggleTask(req *dto.ToggleTaskRequest) (*dto.ToggleTaskResponse, error) {
	tx := s.db.Begin()
	if tx.Error != nil {
		slog.Error("failed to begin transaction", "task_id", req.ID, "error", tx.Error)
		return nil, tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	task, err := s.taskRepo.ToggleTask(tx, req.ID)
	if err != nil {
		tx.Rollback()
		slog.Error("failed to toggle task", "task_id", req.ID, "error", err)
		return nil, err
	}

	if task.IsCompleted {
		if err := s.feedRepo.Create(tx, &model.Feed{
			ActorID: task.UserID,
			TaskID:  task.ID,
			Type:    model.TaskCompleted,
		}); err != nil {
			tx.Rollback()
			slog.Error("failed to create completion feed", "task_id", task.ID, "error", err)
			return nil, err
		}
	} else {
		// 토글 해제 시 완료 피드 삭제
		if err := s.feedRepo.DeleteByActorAndTask(tx, task.UserID, task.ID, model.TaskCompleted); err != nil {
			tx.Rollback()
			slog.Error("failed to delete completion feed", "task_id", task.ID, "error", err)
			return nil, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		slog.Error("failed to commit task toggle", "task_id", req.ID, "error", err)
		return nil, err
	}

	return &dto.ToggleTaskResponse{
		ID:          task.ID,
		IsCompleted: task.IsCompleted,
	}, nil
}
