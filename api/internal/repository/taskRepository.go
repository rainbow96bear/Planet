package repository

import (
	"planet/internal/dto"
	"planet/internal/model"
	"time"

	"gorm.io/gorm"
)

type TaskRepository interface {
	CreateTask(tx *gorm.DB, task *model.Task) error
	DeleteTask(tx *gorm.DB, taskId string) error
	UpdateTask(tx *gorm.DB, task *model.Task) error
	GetTaskByID(taskId string) (*model.Task, error)
	GetTasksByMonth(userid string, year, month int, isOwner bool) ([]*model.Task, error)
	GetOrbitSchedulesByMonth(orbiterID string, year, month int) ([]*dto.OrbitScheduleResponse, error)
	ToggleTask(tx *gorm.DB, taskId string) (*model.Task, error)
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) getDB(tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx
	}
	return r.db
}

func (r *taskRepository) CreateTask(tx *gorm.DB, task *model.Task) error {
	return r.getDB(tx).Create(task).Error
}

// UpdateTask는 Title/Description/StartAt/EndAt/IsPublic만 갱신한다.
// UserID/IsCompleted/DeletedAt 등은 이 경로로 바뀌지 않도록 Select로 컬럼을 한정한다.
func (r *taskRepository) UpdateTask(tx *gorm.DB, task *model.Task) error {
	return r.getDB(tx).Model(&model.Task{}).
		Where("id = ?", task.ID).
		Select("Title", "Description", "StartAt", "EndAt", "IsPublic").
		Updates(task).Error
}

func (r *taskRepository) DeleteTask(tx *gorm.DB, taskId string) error {
	return r.getDB(tx).Where("id = ?", taskId).Delete(&model.Task{}).Error
}

func (r *taskRepository) GetTaskByID(taskId string) (*model.Task, error) {
	var task model.Task
	if err := r.db.Where("id = ?", taskId).First(&task).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *taskRepository) GetTasksByMonth(userid string, year, month int, isOwner bool) ([]*model.Task, error) {
	var tasks []*model.Task

	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0)

	query := r.db.Model(&model.Task{}).
		Where("user_id = ? AND start_at >= ? AND start_at < ?", userid, startDate, endDate)

	if !isOwner {
		query = query.Where("is_public = ?", true)
	}

	if err := query.Find(&tasks).Error; err != nil {
		return nil, err
	}

	return tasks, nil
}

// GetOrbitSchedulesByMonth는 orbiterID(요청자)가 Orbit한 사용자들의
// Public Task만 Calendar 조회 범위(year/month) 내에서 가져온다.
//
// Private Task는 SQL WHERE 절에서 원천적으로 제외된다 — Backend가 Private
// 데이터를 응답에 담아 Frontend로 내려보내는 일 자체가 없다.
//
// feedRepository.FindFeed와 동일하게, orbits 테이블을 서브쿼리로 IN 절에 넣어
// "내가 Orbit한 사람 목록"과 "그 사람들의 Task 조회"를 한 쿼리로 처리한다.
func (r *taskRepository) GetOrbitSchedulesByMonth(orbiterID string, year, month int) ([]*dto.OrbitScheduleResponse, error) {
	var result []*dto.OrbitScheduleResponse

	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0)

	err := r.db.
		Table("tasks t").
		Select(`
			t.id,
			t.user_id,
			u.nickname      AS nickname,
			u.profile_image AS profile_image,
			t.title,
			t.start_at,
			t.end_at
		`).
		Joins("JOIN users u ON u.id = t.user_id AND u.deleted_at IS NULL").
		Where("t.user_id IN (?)",
			r.db.Table("orbits").
				Select("orbited_id").
				Where("orbiter_id = ?", orbiterID),
		).
		Where("t.is_public = ?", true).
		Where("t.deleted_at IS NULL").
		Where("t.start_at >= ? AND t.start_at < ?", startDate, endDate).
		Order("t.start_at ASC").
		Scan(&result).Error

	return result, err
}

func (r *taskRepository) ToggleTask(tx *gorm.DB, taskId string) (*model.Task, error) {
	var task model.Task
	if err := r.getDB(tx).Where("id = ?", taskId).First(&task).Error; err != nil {
		return nil, err
	}
	task.IsCompleted = !task.IsCompleted
	if err := r.getDB(tx).Model(&task).Update("is_completed", task.IsCompleted).Error; err != nil {
		return nil, err
	}
	return &task, nil
}
