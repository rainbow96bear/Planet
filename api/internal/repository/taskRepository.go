package repository

import (
	"gorm.io/gorm"
)

type TaskRepository interface {
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db: db}
}
