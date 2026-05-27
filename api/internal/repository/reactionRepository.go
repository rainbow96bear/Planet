package repository

import (
	"planet/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ReactionRepository interface {
	Upsert(reaction *model.Reaction) error
	Delete(taskID, userID string, reactionType model.ReactionType) error
}

type reactionRepository struct {
	db *gorm.DB
}

func NewReactionRepository(db *gorm.DB) ReactionRepository {
	return &reactionRepository{db: db}
}

func (r *reactionRepository) Upsert(reaction *model.Reaction) error {
	return r.db.
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(reaction).Error
}

func (r *reactionRepository) Delete(taskID, userID string, reactionType model.ReactionType) error {
	return r.db.
		Where("task_id = ? AND user_id = ? AND type = ?", taskID, userID, reactionType).
		Delete(&model.Reaction{}).Error
}
