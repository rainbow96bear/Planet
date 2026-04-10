package repository

import (
	"planet/internal/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(tx *gorm.DB, u *model.User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(tx *gorm.DB, u *model.User) error {
	return tx.Create(u).Error
}
