package repository

import (
	"planet/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ReactionRepository interface {
	Upsert(tx *gorm.DB, reaction *model.Reaction) error
	Delete(tx *gorm.DB, taskID, userID string, reactionType model.ReactionType) error
	DeleteByTaskID(tx *gorm.DB, taskID string) error
}

type reactionRepository struct {
	db *gorm.DB
}

func NewReactionRepository(db *gorm.DB) ReactionRepository {
	return &reactionRepository{db: db}
}

func (r *reactionRepository) getDB(tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx
	}
	return r.db
}

func (r *reactionRepository) Upsert(tx *gorm.DB, reaction *model.Reaction) error {
	return r.getDB(tx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(reaction).Error
}

func (r *reactionRepository) Delete(tx *gorm.DB, taskID, userID string, reactionType model.ReactionType) error {
	return r.getDB(tx).
		Where("task_id = ? AND user_id = ? AND type = ?", taskID, userID, reactionType).
		Delete(&model.Reaction{}).Error
}

func (r *reactionRepository) DeleteByTaskID(tx *gorm.DB, taskID string) error {
	return r.getDB(tx).
		Where("task_id = ?", taskID).
		Delete(&model.Reaction{}).Error
}
