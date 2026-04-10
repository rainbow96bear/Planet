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
	"planet/internal/repository"
	"planet/internal/service"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Auto migrate
	if err := db.AutoMigrate(&model.User{}); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	userRepo := repository.NewUserRepository(db)

	userSvc := service.NewUserService(db, userRepo)

	userHandler := handler.NewUserHandler(userSvc)

	r := gin.Default()
	handler.RegisterRoutes(r, userHandler)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.App.Port),
		Handler: r,
	}

	go func() {
		log.Printf("server runnig on : %s\n", cfg.App.Port)
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
