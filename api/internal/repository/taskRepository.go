package repository

import (
	"planet/internal/model"

	"gorm.io/gorm"
)

type TaskRepository interface {
	CreateTask(tx *gorm.DB, task *model.Task) error
	DeleteTask(tx *gorm.DB, taskId uint) error
	GetTaskByID(taskId uint) (*model.Task, error)
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) CreateTask(tx *gorm.DB, task *model.Task) error {
	return tx.Create(task).Error
}

func (r *taskRepository) DeleteTask(tx *gorm.DB, taskId uint) error {
	return tx.Delete(&model.Task{}, taskId).Error
}

func (r *taskRepository) GetTaskByID(taskId uint) (*model.Task, error) {
	var task model.Task
	if err := r.db.First(&task, taskId).Error; err != nil {
		return nil, err
	}
	return &task, nil
}
