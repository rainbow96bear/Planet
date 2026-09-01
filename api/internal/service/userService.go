package service

import (
	"errors"
	"log/slog"
	"planet/internal/dto"
	"planet/internal/model"
	"planet/internal/repository"
	"planet/internal/storage"

	"gorm.io/gorm"
)

type UserService interface {
	GetProfile(*dto.GetProfileRequest) (*dto.GetProfileResponse, error)
	EnterOrbit(*dto.EnterOrbitRequest) (*dto.EnterOrbitResponse, error)
	LeaveOrbit(*dto.LeaveOrbitRequest) (*dto.LeaveOrbitResponse, error)
	UpdateProfile(*dto.UpdateProfileRequest) (*dto.UpdateProfileResponse, error)
	UploadProfileImage(*dto.UploadProfileImageRequest) (*dto.UploadProfileImageResponse, error)
	DeleteProfileImage(userID string) error
}

type userService struct {
	db               *gorm.DB
	userRepo         repository.UserRepository
	orbitRepo        repository.OrbitRepository
	taskRepo         repository.TaskRepository
	feedRepo         repository.FeedRepository
	notificationRepo repository.NotificationRepository
	fileStorage      storage.FileStorage
}

func NewUserService(
	db *gorm.DB,
	userRepo repository.UserRepository,
	orbitRepo repository.OrbitRepository,
	taskRepo repository.TaskRepository,
	feedRepo repository.FeedRepository,
	notificationRepo repository.NotificationRepository,
	fileStorage storage.FileStorage,

) UserService {
	return &userService{
		db:               db,
		userRepo:         userRepo,
		orbitRepo:        orbitRepo,
		taskRepo:         taskRepo,
		feedRepo:         feedRepo,
		notificationRepo: notificationRepo,
		fileStorage:      fileStorage,
	}
}

// service
func (s *userService) GetProfile(req *dto.GetProfileRequest) (*dto.GetProfileResponse, error) {
	user, err := s.userRepo.FindByUserId(req.UserId)
	if err != nil {
		return nil, err
	}

	gravity, err := s.orbitRepo.CountGravity(req.UserId)
	if err != nil {
		return nil, err
	}

	orbit, err := s.orbitRepo.CountOrbit(req.UserId)
	if err != nil {
		return nil, err
	}

	isOwner := req.UserId == req.RequesterUserId

	var isOrbiting bool
	if !isOwner {
		isOrbiting, err = s.orbitRepo.IsOrbiting(req.RequesterUserId, req.UserId)
		if err != nil {
			return nil, err
		}
	}

	return &dto.GetProfileResponse{
		Username:     user.Username,
		Nickname:     user.Nickname,
		ProfileImage: user.ProfileImage,
		IsOwner:      isOwner,
		IsOrbiting:   isOrbiting,
		Gravity:      gravity,
		Orbit:        orbit,
	}, nil
}

func (s *userService) EnterOrbit(req *dto.EnterOrbitRequest) (*dto.EnterOrbitResponse, error) {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r) // 롤백 후 패닉을 다시 던져서 상위(recovery 미들웨어)가 처리하도록 함
		}
	}()

	if err := s.orbitRepo.EnterOrbit(tx, &model.Orbit{
		OrbiterID: req.OrbiterID,
		OrbitedID: req.OrbitedID,
	}); err != nil {
		tx.Rollback()
		return nil, errors.New("이미 궤도에 진입했습니다")
	}

	if err := s.notificationRepo.Upsert(tx, &model.Notification{
		ReceiverID: req.OrbitedID,
		ActorID:    req.OrbiterID,
		TargetID:   req.OrbiterID,
		TargetType: model.NotificationTargetTypeUser,
		Type:       model.NotificationTypeOrbitEntered,
	}); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &dto.EnterOrbitResponse{IsOrbiting: true}, nil
}

func (s *userService) LeaveOrbit(req *dto.LeaveOrbitRequest) (*dto.LeaveOrbitResponse, error) {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if err := s.orbitRepo.LeaveOrbit(tx, req.OrbiterID, req.OrbitedID); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := s.notificationRepo.DeleteByActorAndTask(
		tx,
		req.OrbiterID,
		req.OrbiterID,
		model.NotificationTypeOrbitEntered,
	); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &dto.LeaveOrbitResponse{IsOrbiting: false}, nil
}

func (s *userService) UpdateProfile(req *dto.UpdateProfileRequest) (*dto.UpdateProfileResponse, error) {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	user := &model.User{
		BaseModel: model.BaseModel{ID: req.UserID},
		Nickname:  req.Nickname,
	}

	if err := s.userRepo.UpdateProfile(tx, user); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &dto.UpdateProfileResponse{
		UserID:   user.ID,
		Username: user.Username,
		Nickname: user.Nickname,
	}, nil
}

func (s *userService) UploadProfileImage(req *dto.UploadProfileImageRequest) (*dto.UploadProfileImageResponse, error) {
	// 기존 이미지 URL을 먼저 조회 (업로드 성공 후 이전 파일 정리용)
	existing, err := s.userRepo.FindByUserId(req.UserID)
	if err != nil {
		return nil, err
	}
	previousImageURL := existing.ProfileImage

	// 1) 스토리지 업로드는 트랜잭션 밖에서 먼저 수행 (DB 트랜잭션 안에서 외부 I/O를 잡아두지 않기 위함)
	newURL, err := s.fileStorage.Upload(storage.UploadInput{
		Key:         storage.ProfileImageKey(req.UserID, req.Filename),
		Reader:      req.File,
		ContentType: req.ContentType,
		Size:        req.Size,
	})
	if err != nil {
		return nil, errors.New("이미지 업로드에 실패했습니다")
	}

	// 2) DB 반영은 트랜잭션으로 처리. UpdateProfileImage는 map 기반이라 빈 문자열도 정확히 반영된다.
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if err := s.userRepo.UpdateProfileImage(tx, req.UserID, newURL); err != nil {
		tx.Rollback()
		// DB 반영 실패 시 방금 업로드한 파일을 정리 (베스트 에포트)
		if delErr := s.fileStorage.Delete(newURL); delErr != nil {
			slog.Error("failed to cleanup orphaned profile image after db failure", "url", newURL, "err", delErr)
		}
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		if delErr := s.fileStorage.Delete(newURL); delErr != nil {
			slog.Error("failed to cleanup orphaned profile image after commit failure", "url", newURL, "err", delErr)
		}
		return nil, err
	}

	// 3) 이전 이미지가 있었다면 커밋 성공 후 정리 (실패해도 응답에는 영향 없음 — 베스트 에포트)
	if previousImageURL != "" && previousImageURL != newURL {
		if delErr := s.fileStorage.Delete(previousImageURL); delErr != nil {
			slog.Error("failed to delete previous profile image", "url", previousImageURL, "err", delErr)
		}
	}

	return &dto.UploadProfileImageResponse{ProfileImage: newURL}, nil
}

func (s *userService) DeleteProfileImage(userID string) error {
	existing, err := s.userRepo.FindByUserId(userID)
	if err != nil {
		return err
	}
	if existing.ProfileImage == "" {
		return nil // 이미 기본 이미지 상태 — 별도 처리 불필요
	}

	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if err := s.userRepo.UpdateProfileImage(tx, userID, ""); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	if delErr := s.fileStorage.Delete(existing.ProfileImage); delErr != nil {
		slog.Error("failed to delete profile image file", "url", existing.ProfileImage, "err", delErr)
	}

	return nil
}
