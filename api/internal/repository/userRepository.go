package repository

import (
	"planet/internal/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateWithTx(u *model.User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateWithTx(u *model.User) error {
	tx := r.db.Begin()

	if err := tx.Create(u).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
