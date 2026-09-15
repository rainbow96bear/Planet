package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"planet/internal/config"
	"planet/internal/database"
	"planet/internal/handler"
	"planet/internal/model"
	"planet/internal/pkg"
	"planet/internal/pkg/logger"
	"planet/internal/repository"
	"planet/internal/service"
	"planet/internal/storage"
	"runtime/debug"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	appLogger := logger.New()
	slog.SetDefault(appLogger)

	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	pkg.InitToken(cfg)

	db, err := database.Connect(cfg)
	if err != nil {
		slog.Error("failed to connect database", "error", err)
		os.Exit(1)
	}
	slog.Info("database connected")

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
			slog.Info("migration skipped (already exists)")
		} else {
			slog.Error("migration error", "error", err)
			os.Exit(1)
		}
	} else {
		slog.Info("migration completed")
	}

	// 파일 스토리지 — 현재는 로컬 디스크 구현체만 존재.
	// 배포 플랫폼이 정해지면 이 부분만 다른 구현체로 교체하면 됨 (서비스/핸들러 코드는 변경 불필요).
	uploadDir := "./uploads"
	baseURL := fmt.Sprintf("http://localhost:%s/uploads", cfg.App.Port)
	if err := os.MkdirAll(uploadDir, 0o755); err != nil {
		slog.Error("failed to create upload dir", "dir", uploadDir, "error", err)
		os.Exit(1)
	}
	fileStorage := storage.NewLocalFileStorage(uploadDir, baseURL)
	slog.Info("file storage initialized", "backend", "local", "upload_dir", uploadDir)

	userRepo := repository.NewUserRepository(db)
	taskRepo := repository.NewTaskRepository(db)
	orbitRepo := repository.NewOrbitRepository(db)
	feedRepo := repository.NewFeedRepository(db)
	notificationRepo := repository.NewNotificationRepository(db)
	reactionRepo := repository.NewReactionRepository(db)

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

	// gin.Default()의 텍스트 포맷 Logger 미들웨어 대신 slog 기반 요청 로거를 사용해서
	// 서버 전체 로그 포맷을 일관되게 유지한다 (Cloud Run 등에서 로그 수집 시 유리).
	//
	// 미들웨어 순서 주의: slogRequestLogger가 바깥쪽, slogRecovery가 안쪽이어야 한다.
	// gin.Recovery()는 쓰지 않는다 — panic이 나면 slogRequestLogger의 c.Next() 이후
	// 코드(status/latency 계산, 로그 출력)가 실행되지 않고 건너뛰어져서, JSON 로그에
	// 500 에러가 전혀 남지 않는 문제가 있었다. slogRecovery가 panic을 흡수해서
	// 정상적으로 500을 반환해야 바깥쪽 로거가 그 status를 로깅할 수 있다.
	r := gin.New()
	r.Use(slogRequestLogger())
	r.Use(slogRecovery())

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
		slog.Info("server running", "port", cfg.App.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("forced shutdown", "error", err)
		os.Exit(1)
	}
	slog.Info("server exited")
}

// slogRequestLogger는 각 HTTP 요청을 slog로 구조화해서 남긴다.
// 4xx/5xx 응답은 Warn/Error로, 그 외에는 Info로 구분한다.
// user_id(인증 미들웨어가 설정한 경우)와 핸들러 에러(c.Errors)가 있으면 함께 로깅한다.
func slogRequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		c.Next()

		status := c.Writer.Status()
		attrs := []any{
			"method", c.Request.Method,
			"path", path,
			"status", status,
			"latency_ms", time.Since(start).Milliseconds(),
			"client_ip", c.ClientIP(),
		}

		if uid, exists := c.Get("userID"); exists {
			attrs = append(attrs, "user_id", uid)
		}

		if len(c.Errors) > 0 {
			attrs = append(attrs, "error", c.Errors.String())
		}

		switch {
		case status >= 500:
			slog.Error("http request", attrs...)
		case status >= 400:
			slog.Warn("http request", attrs...)
		default:
			slog.Info("http request", attrs...)
		}
	}
}

// slogRecovery는 gin.Recovery()를 대체한다. panic을 잡아 slog로 스택트레이스와 함께
// 남기고 500을 반환한다. slogRequestLogger보다 안쪽에 등록해야, panic이 나도
// 바깥쪽 로거의 status/latency 로깅이 정상적으로 실행된다.
func slogRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				slog.Error("panic recovered",
					"method", c.Request.Method,
					"path", c.Request.URL.Path,
					"error", fmt.Sprint(rec),
					"stack", string(debug.Stack()),
				)
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
