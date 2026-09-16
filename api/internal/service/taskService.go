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
	GetOrbitSchedulesByMonth(req *dto.GetOrbitSchedulesByMonthRequest) ([]*dto.OrbitScheduleResponse, error)
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
	task := &model.Task{
		UserID:      req.UserID,
		Title:       req.Title,
		Description: req.Description,
		StartAt:     req.StartAt,
		EndAt:       req.EndAt,
		IsPublic:    req.IsPublic,
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.taskRepo.CreateTask(tx, task); err != nil {
			slog.Error("failed to create task", "user_id", req.UserID, "error", err)
			return err
		}

		// TargetID/TargetType → TaskID, UserID → ActorID
		if err := s.feedRepo.Create(tx, &model.Feed{
			ActorID: req.UserID,
			TaskID:  task.ID,
			Type:    model.TaskCreated,
		}); err != nil {
			slog.Error("failed to create feed for task", "task_id", task.ID, "user_id", req.UserID, "error", err)
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	slog.Info("task created", "task_id", task.ID, "user_id", req.UserID, "is_public", task.IsPublic)

	return &dto.CreateTaskResponse{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		StartAt:     task.StartAt,
		EndAt:       task.EndAt,
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

	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.taskRepo.DeleteTask(tx, req.ID); err != nil {
			slog.Error("failed to delete task", "task_id", req.ID, "user_id", req.UserID, "error", err)
			return err
		}

		if err := s.feedRepo.DeleteByTaskID(tx, req.ID); err != nil {
			slog.Error("failed to delete feed for task", "task_id", req.ID, "error", err)
			return err
		}

		if err := s.reactionRepo.DeleteByTaskID(tx, req.ID); err != nil {
			slog.Error("failed to delete reactions for task", "task_id", req.ID, "error", err)
			return err
		}

		return nil
	})
	if err != nil {
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
			StartAt:     task.StartAt,
			EndAt:       task.EndAt,
			IsCompleted: task.IsCompleted,
			IsPublic:    task.IsPublic,
		}
	}
	return result, nil
}

func (s *taskService) GetOrbitSchedulesByMonth(req *dto.GetOrbitSchedulesByMonthRequest) ([]*dto.OrbitScheduleResponse, error) {
	schedules, err := s.taskRepo.GetOrbitSchedulesByMonth(req.OrbiterID, req.Year, req.Month)
	if err != nil {
		slog.Error("failed to fetch orbit schedules",
			"orbiter_id", req.OrbiterID,
			"year", req.Year,
			"month", req.Month,
			"error", err,
		)
		return nil, err
	}
	return schedules, nil
}

func (s *taskService) ToggleTask(req *dto.ToggleTaskRequest) (*dto.ToggleTaskResponse, error) {
	existing, err := s.taskRepo.GetTaskByID(req.ID)
	if err != nil {
		return nil, errors.New("존재하지 않는 할 일입니다")
	}
	if existing.UserID != req.UserID {
		slog.Warn("unauthorized task toggle attempt",
			"task_id", req.ID,
			"owner_id", existing.UserID,
			"requester_id", req.UserID,
		)
		return nil, errors.New("권한이 없습니다")
	}

	var task *model.Task
	err = s.db.Transaction(func(tx *gorm.DB) error {
		toggled, err := s.taskRepo.ToggleTask(tx, req.ID)
		if err != nil {
			slog.Error("failed to toggle task", "task_id", req.ID, "error", err)
			return err
		}
		task = toggled

		if task.IsCompleted {
			if err := s.feedRepo.Create(tx, &model.Feed{
				ActorID: task.UserID,
				TaskID:  task.ID,
				Type:    model.TaskCompleted,
			}); err != nil {
				slog.Error("failed to create completion feed", "task_id", task.ID, "error", err)
				return err
			}
		} else {
			if err := s.feedRepo.DeleteByActorAndTask(tx, task.UserID, task.ID, model.TaskCompleted); err != nil {
				slog.Error("failed to delete completion feed", "task_id", task.ID, "error", err)
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &dto.ToggleTaskResponse{
		ID:          task.ID,
		IsCompleted: task.IsCompleted,
	}, nil
}
