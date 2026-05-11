package repository

import (
	"planet/internal/model"
	"time"

	"gorm.io/gorm"
)

type TaskRepository interface {
	CreateTask(tx *gorm.DB, task *model.Task) error
	DeleteTask(tx *gorm.DB, taskId string) error
	GetTaskByID(taskId string) (*model.Task, error)
	GetTasksByMonth(userid string, year, month int, isOwner bool) ([]*model.Task, error)
	ToggleTask(tx *gorm.DB, taskId string) (*model.Task, error)
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

func (r *taskRepository) DeleteTask(tx *gorm.DB, taskId string) error {
	return tx.Delete(&model.Task{}, taskId).Error
}

func (r *taskRepository) GetTaskByID(taskId string) (*model.Task, error) {
	var task model.Task
	if err := r.db.First(&task, taskId).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *taskRepository) GetTasksByMonth(userid string, year, month int, isOwner bool) ([]*model.Task, error) {
	var tasks []*model.Task

	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0)
	query := r.db.Model(&model.Task{}).
		Select("tasks.*").
		Joins("JOIN users ON users.id = tasks.user_id").
		Where("users.id = ? AND tasks.date >= ? AND tasks.date < ?", userid, startDate, endDate)

	if !isOwner {
		query = query.Where("tasks.is_public = ?", true)
	}

	if err := query.Find(&tasks).Error; err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *taskRepository) ToggleTask(tx *gorm.DB, taskId string) (*model.Task, error) {
	var task model.Task
	if err := tx.First(&task, taskId).Error; err != nil {
		return nil, err
	}
	task.IsCompleted = !task.IsCompleted
	if err := tx.Model(&task).Update("is_completed", task.IsCompleted).Error; err != nil {
		return nil, err
	}
	return &task, nil
}
