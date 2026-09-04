package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"planet/internal/config"
	"planet/internal/database"
	"planet/internal/handler"
	"planet/internal/model"
	"planet/internal/pkg"
	"planet/internal/repository"
	"planet/internal/service"
	"planet/internal/storage"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	pkg.InitToken(cfg)

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Auto migrate
	if err := db.AutoMigrate(
		&model.User{},
		&model.Task{},
		&model.Orbit{},
		&model.Feed{},
		&model.Notification{},
		&model.Reaction{},
	); err != nil {
		if strings.Contains(err.Error(), "already exists") {
			log.Println("migration skipped (already exists)")
		} else {
			log.Fatalf("migration error: %v", err)
		}
	}

	// 파일 스토리지 — 현재는 로컬 디스크 구현체만 존재.
	// 배포 플랫폼이 정해지면 이 부분만 다른 구현체로 교체하면 됨 (서비스/핸들러 코드는 변경 불필요).
	uploadDir := "./uploads"
	baseURL := fmt.Sprintf("http://localhost:%s/uploads", cfg.App.Port)
	if err := os.MkdirAll(uploadDir, 0o755); err != nil {
		log.Fatalf("failed to create upload dir: %v", err)
	}
	fileStorage := storage.NewLocalFileStorage(uploadDir, baseURL)

	userRepo := repository.NewUserRepository(db)
	taskRepo := repository.NewTaskRepository(db)
	orbitRepo := repository.NewOrbitRepository(db)
	feedRepo := repository.NewFeedRepository(db)
	notificationRepo := repository.NewNotificationRepository(db)
	reactionRepo := repository.NewReactionRepository(db)

	authSvc := service.NewAuthService(db, userRepo, fileStorage)
	authSvc := service.NewAuthService(db, userRepo, fileStorage)
	taskSvc := service.NewTaskService(db, taskRepo, feedRepo, reactionRepo)
	userSvc := service.NewUserService(db, userRepo, orbitRepo, taskRepo, feedRepo, notificationRepo, fileStorage)
	searchSvc := service.NewSearchService(db, userRepo, orbitRepo)
	feedSvc := service.NewFeedService(db, feedRepo)
	notificationSvc := service.NewNotificationService(db, notificationRepo, userRepo)
	reactionSvc := service.NewReactionService(db, reactionRepo, taskRepo, notificationRepo)

	authHandler := handler.NewAuthHandler(authSvc)
	taskHandler := handler.NewTaskHandler(taskSvc)
	userHandler := handler.NewUserHandler(userSvc)
	searchHandler := handler.NewSearchHandler(searchSvc)
	feedHandler := handler.NewFeedHandler(feedSvc)
	notificationHandler := handler.NewNotificationHandler(notificationSvc)
	reactionHandler := handler.NewReactionHandler(reactionSvc)

	r := gin.Default()

	// 로컬 스토리지에 저장된 파일 서빙 (프로덕션 스토리지로 교체되면 이 라인은 제거)
	r.Static("/uploads", uploadDir)

	handler.RegisterRoutes(
		r,
		authHandler,
		taskHandler,
		userHandler,
		searchHandler,
		feedHandler,
		notificationHandler,
		reactionHandler,
	)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.App.Port),
		Handler: r,
	}

	go func() {
		log.Printf("server running on : %s\n", cfg.App.Port)
		log.Printf("server running on : %s\n", cfg.App.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Printf("shutting down...\n")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("forced shutdown: %v", err)
	}
	log.Printf("server exited\n")
}
