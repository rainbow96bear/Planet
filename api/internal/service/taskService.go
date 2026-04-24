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
}

type taskService struct {
	db *gorm.DB
	// userRepo repository.UserRepository
	taskRepo repository.TaskRepository
}

func NewTaskService(db *gorm.DB, taskRepo repository.TaskRepository) TaskService {
	return &taskService{
		db:       db,
		taskRepo: taskRepo,
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
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &dto.CreateTaskResponse{
		ID:       task.ID,
		Title:    task.Title,
		Date:     task.Date,
		IsPublic: task.IsPublic,
	}, nil
}

func (s *taskService) DeleteTask(req *dto.DeleteTaskRequest) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	task, err := s.taskRepo.GetTaskByID(req.TaskID)
	if err != nil {
		return errors.New("존재하지 않는 할 일입니다")
	}

	if task.UserID != req.UserID {
		return errors.New("권한이 없습니다")
	}

	if err := s.taskRepo.DeleteTask(tx, req.TaskID); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
