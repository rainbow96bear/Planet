package service

import (
	"errors"
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

		return nil, err
	}

	// TargetID/TargetType → TaskID, UserID → ActorID
	if err := s.feedRepo.Create(tx, &model.Feed{
		ActorID: req.UserID,
		TaskID:  task.ID,
		Type:    model.TaskCreated,
	}); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

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
		return errors.New("권한이 없습니다")
	}

	tx := s.db.Begin()

	if err := s.taskRepo.DeleteTask(tx, req.ID); err != nil {
		tx.Rollback()
		return err
	}

	if err := s.feedRepo.DeleteByTaskID(tx, req.ID); err != nil {
		tx.Rollback()
		return err
	}

	if err := s.reactionRepo.DeleteByTaskID(tx, req.ID); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *taskService) GetTasksByMonth(req *dto.GetTasksByMonthRequest) ([]*dto.GetTasksByMonthResponse, error) {
	isOwner := req.UserID == req.RequesterUserId

	tasks, err := s.taskRepo.GetTasksByMonth(req.UserID, req.Year, req.Month, isOwner)
	if err != nil {
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
	return result, nil // fmt.Printf 제거
}

func (s *taskService) ToggleTask(req *dto.ToggleTaskRequest) (*dto.ToggleTaskResponse, error) {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	task, err := s.taskRepo.ToggleTask(tx, req.ID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if task.IsCompleted {
		if err := s.feedRepo.Create(tx, &model.Feed{
			ActorID: task.UserID,
			TaskID:  task.ID,
			Type:    model.TaskCompleted,
		}); err != nil {
			tx.Rollback()
			return nil, err
		}
	} else {
		// 토글 해제 시 완료 피드 삭제
		if err := s.feedRepo.DeleteByActorAndTask(tx, task.UserID, task.ID, model.TaskCompleted); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &dto.ToggleTaskResponse{
		ID:          task.ID,
		IsCompleted: task.IsCompleted,
	}, nil
}
