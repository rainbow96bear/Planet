package service

import (
	"gorm.io/gorm"
)

type TaskService interface {
}

type taskService struct {
	db *gorm.DB
	// userRepo repository.UserRepository
}

func NewTaskService(db *gorm.DB) TaskService {
	return &taskService{db: db}
}
