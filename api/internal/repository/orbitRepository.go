package repository

import (
	"errors"
	"planet/internal/model"

	"gorm.io/gorm"
)

type OrbitRepository interface {
	EnterOrbit(tx *gorm.DB, o *model.Orbit) error
	LeaveOrbit(tx *gorm.DB, orbiterID, orbitedID string) error
	IsOrbiting(orbiterID, orbitedID string) (bool, error)
	CountGravity(userID string) (int64, error)
	CountOrbit(userID string) (int64, error)
}

type orbitRepository struct {
	db *gorm.DB
}

func NewOrbitRepository(db *gorm.DB) OrbitRepository {
	return &orbitRepository{db: db}
}

func (r *orbitRepository) getDB(tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx
	}
	return r.db
}

func (r *orbitRepository) EnterOrbit(tx *gorm.DB, o *model.Orbit) error {
	return r.getDB(tx).Create(o).Error
}

func (r *orbitRepository) LeaveOrbit(tx *gorm.DB, orbiterID, orbitedID string) error {
	return r.getDB(tx).Where("orbiter_id = ? AND orbited_id = ?", orbiterID, orbitedID).Delete(&model.Orbit{}).Error
}

func (r *orbitRepository) IsOrbiting(orbiterID, orbitedID string) (bool, error) {
	var orbit model.Orbit
	err := r.db.
		Where("orbiter_id = ? AND orbited_id = ?", orbiterID, orbitedID).
		Limit(1).
		First(&orbit).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return err == nil, err
}

// CountGravity는 userID에게 이끌려온(=팔로우하는) 사람 수를 센다. (기존 CountFollowers)
func (r *orbitRepository) CountGravity(userID string) (int64, error) {
	var count int64
	err := r.db.Model(&model.Orbit{}).Where("orbited_id = ?", userID).Count(&count).Error
	return count, err
}

// CountOrbit은 userID가 궤도를 돌고 있는(=팔로우하는) 대상 수를 센다. (기존 CountFollowing)
func (r *orbitRepository) CountOrbit(userID string) (int64, error) {
	var count int64
	err := r.db.Model(&model.Orbit{}).Where("orbiter_id = ?", userID).Count(&count).Error
	return count, err
}
