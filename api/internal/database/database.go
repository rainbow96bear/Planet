package database

import (
	"log"
	"planet/internal/config"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(cfg *config.Config) (*gorm.DB, error) {
	log.Printf("DB: %s", cfg.DB)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  cfg.DB.DSN(),
		PreferSimpleProtocol: true, // PgBouncer/Supavisor 풀링 모드에서 prepared statement 충돌 방지
	}), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)

	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	log.Println("database connected")

	return db, nil
}
