package service

import (
	"errors"
	"planet/internal/dto"
	"planet/internal/model"
	"planet/internal/repository"

	"gorm.io/gorm"
)

type UserService interface {
	GetProfile(*dto.GetProfileRequest) (*dto.GetProfileResponse, error)
	Follow(*dto.FollowRequest) (*dto.FollowResponse, error)
	Unfollow(*dto.UnfollowRequest) (*dto.UnfollowResponse, error)
	UpdateProfile(*dto.UpdateProfileRequest) (*dto.UpdateProfileResponse, error)
	GetActivity(*dto.GetActivityRequest) (*dto.GetActivityResponse, error)
}

type userService struct {
	db               *gorm.DB
	userRepo         repository.UserRepository
	followRepo       repository.FollowRepository
	taskRepo         repository.TaskRepository
	activityRepo     repository.ActivityRepository
	notificationRepo repository.NotificationRepository
}

func NewUserService(
	db *gorm.DB,
	userRepo repository.UserRepository,
	followRepo repository.FollowRepository,
	taskRepo repository.TaskRepository,
	activityRepo repository.ActivityRepository,
	notificationRepo repository.NotificationRepository,

) UserService {
	return &userService{
		db:               db,
		userRepo:         userRepo,
		followRepo:       followRepo,
		taskRepo:         taskRepo,
		activityRepo:     activityRepo,
		notificationRepo: notificationRepo,
	}
}

// service
func (s *userService) GetProfile(req *dto.GetProfileRequest) (*dto.GetProfileResponse, error) {
	user, err := s.userRepo.FindByUserId(req.UserId)
	if err != nil {
		return nil, err
	}

	followers, err := s.followRepo.CountFollowers(req.UserId)
	if err != nil {
		return nil, err
	}

	following, err := s.followRepo.CountFollowing(req.UserId)
	if err != nil {
		return nil, err
	}

	isOwner := req.UserId == req.RequesterUserId

	var isFollowing bool
	if !isOwner {
		isFollowing, err = s.followRepo.IsFollowing(req.RequesterUserId, req.UserId)
		if err != nil {
			return nil, err
		}
	}

	return &dto.GetProfileResponse{
		Username:    user.Username,
		Nickname:    user.Nickname,
		IsOwner:     req.UserId == req.RequesterUserId,
		IsFollowing: isFollowing,
		Followers:   followers,
		Following:   following,
	}, nil
}

func (s *userService) Follow(req *dto.FollowRequest) (*dto.FollowResponse, error) {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := s.followRepo.Follow(tx, &model.Follow{
		FollowerID:  req.FollowerID,
		FollowingID: req.FollowingID,
	}); err != nil {
		tx.Rollback()
		return nil, errors.New("이미 팔로우 중입니다")
	}

	// 추가
	// targetType := model.TargetTypeUser
	// if err := s.activityRepo.Create(tx, &model.Activity{
	// 	UserID:     req.FollowerID,
	// 	Type:       model.ActivityFollowed,
	// 	TargetID:   req.FollowingID,
	// 	TargetType: targetType,
	// }); err != nil {
	// 	tx.Rollback()
	// 	return nil, err
	// }

	if err := s.notificationRepo.Upsert(tx, &model.Notification{
		ReceiverID: req.FollowingID,
		ActorID:    req.FollowerID,
		TargetID:   req.FollowerID,
		TargetType: model.NotificationTargetTypeUser,
		Type:       model.NotificationFollowed,
	}); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &dto.FollowResponse{IsFollowing: true}, nil
}

func (s *userService) Unfollow(req *dto.UnfollowRequest) (*dto.UnfollowResponse, error) {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := s.followRepo.Unfollow(tx, req.FollowerID, req.FollowingID); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := s.notificationRepo.DeleteByActorAndTask(
		tx,
		req.FollowerID,
		"", // 팔로우 알림은 TaskID 없으니 빈 문자열
		model.NotificationFollowed,
	); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &dto.UnfollowResponse{IsFollowing: false}, nil
}

func (s *userService) UpdateProfile(req *dto.UpdateProfileRequest) (*dto.UpdateProfileResponse, error) {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
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

func (s *userService) GetActivity(req *dto.GetActivityRequest) (*dto.GetActivityResponse, error) {
	isOwner := req.UserID == req.RequesterUserId

	activity, err := s.activityRepo.FindLatestByUserID(req.UserID)
	if err != nil || activity == nil {
		return nil, err
	}

	switch activity.TargetType {
	case model.TargetTypeTask:
		task, err := s.taskRepo.GetTaskByID(activity.TargetID)
		if err != nil {
			return nil, err
		}
		if !isOwner && !task.IsPublic {
			return nil, nil
		}
		return &dto.GetActivityResponse{
			Type:      string(activity.Type),
			TaskTitle: &task.Title,
			CreatedAt: activity.CreatedAt,
		}, nil

	case model.TargetTypeUser:
		return &dto.GetActivityResponse{
			Type:      string(activity.Type),
			CreatedAt: activity.CreatedAt,
		}, nil
	}

	return nil, nil
}
