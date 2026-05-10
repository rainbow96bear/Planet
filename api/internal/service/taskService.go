package service

import (
	"errors"
	"fmt"
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

	task, err := s.taskRepo.GetTaskByID(req.ID)
	if err != nil {
		return errors.New("존재하지 않는 할 일입니다")
	}

	if task.UserID != req.UserID {
		return errors.New("권한이 없습니다")
	}

	if err := s.taskRepo.DeleteTask(tx, req.ID); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
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
	fmt.Printf("tasks : %+v\n", result)
	return result, nil
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

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &dto.ToggleTaskResponse{
		ID:          task.ID,
		IsCompleted: task.IsCompleted,
	}, nil
}
