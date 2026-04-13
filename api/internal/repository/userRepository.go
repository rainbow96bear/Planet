package repository

import (
	"planet/internal/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(tx *gorm.DB, u *model.User) error
	IsUsernameExists(username string) (bool, error)
	FindByUsername(username string) (model.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(tx *gorm.DB, u *model.User) error {
	return tx.Create(u).Error
}

func (r *userRepository) IsUsernameExists(username string) (bool, error) {
	var count int64
	err := r.db.Model(&model.User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}

func (r *userRepository) FindByUsername(username string) (model.User, error) {
	var user model.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return model.User{}, err
	}
	return user, nil
}
