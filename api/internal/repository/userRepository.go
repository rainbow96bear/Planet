package repository

import (
	"planet/internal/model"
	"strings"
	"time"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(tx *gorm.DB, u *model.User) error
	IsUsernameExists(username string) (bool, error)
	FindByUsername(username string) (model.User, error)
	FindByProviderInfo(provider, providerID string) (*model.User, error)
	FindByUserId(userid string) (model.User, error)
	UpdateProfile(tx *gorm.DB, u *model.User) error
	UpdateProfileImage(tx *gorm.DB, userID string, avatarURL string) error
	SearchByKeyword(q string) ([]*model.User, error)
	UpdateLastLogin(userID string, loginAt time.Time) error
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

func (r *userRepository) FindByProviderInfo(provider, providerID string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("provider = ? AND provider_id = ?", provider, providerID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByUserId(userid string) (model.User, error) {
	var user model.User
	if err := r.db.Where("id = ?", userid).First(&user).Error; err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (r *userRepository) UpdateProfile(tx *gorm.DB, u *model.User) error {
	return tx.Model(u).Updates(model.User{
		Nickname: u.Nickname,
	}).Error
}

// UpdateProfileImage는 map을 사용해 명시적으로 필드를 지정한다.
// 구조체 기반 Updates는 zero value(빈 문자열 등)를 무시하므로,
// 프로필 이미지 삭제(빈 문자열로 되돌리기) 시 반드시 map 또는 Select를 써야 한다.
func (r *userRepository) UpdateProfileImage(tx *gorm.DB, userID string, profileImageURL string) error {
	return tx.Model(&model.User{}).
		Where("id = ?", userID).
		Update("profile_image", profileImageURL).
		Error
}

func (r *userRepository) UpdateLastLogin(userID string, loginAt time.Time) error {
	return r.db.
		Model(&model.User{}).
		Where("id = ?", userID).
		Update("last_login_at", loginAt).
		Error
}

func (r *userRepository) SearchByKeyword(q string) ([]*model.User, error) {
	var users []*model.User
	keyword := "%" + escapeLikeWildcards(q) + "%"
	if err := r.db.
		Where("username LIKE ? OR nickname LIKE ?", keyword, keyword).
		Limit(20).
		Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

// escapeLikeWildcards는 LIKE 패턴에서 특수 취급되는 %, _ 문자를 이스케이프한다.
// 사용자가 검색어에 %나 _를 입력했을 때 의도치 않은 와일드카드로 해석되는 것을 방지한다.
func escapeLikeWildcards(s string) string {
	replacer := strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)
	return replacer.Replace(s)
}
