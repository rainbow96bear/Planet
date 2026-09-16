package database

import (
	"log/slog"

	"gorm.io/gorm"
)

// MigrateTaskScheduleColumns는 기존 Task.Date(날짜만)를
// StartAt/EndAt(시간 포함) 컬럼으로 전환하는 1회성 백필 마이그레이션이다.
//
// 정책: 기존 row는 하루 종일 일정으로 간주한다.
//
//	StartAt = date의 00:00:00
//	EndAt   = date의 23:59:59
//
// IF NOT EXISTS / WHERE start_at IS NULL 가드가 있어 여러 번 실행돼도 안전하다(idempotent).
// db.AutoMigrate(&model.Task{}, ...)보다 반드시 먼저 호출해야 한다 —
// AutoMigrate가 StartAt/EndAt을 NOT NULL로 선언하기 전에 기존 row를 채워둬야 하기 때문이다.
func MigrateTaskScheduleColumns(db *gorm.DB) error {
	// 1) 컬럼을 nullable 상태로 먼저 추가 (NOT NULL은 아직 걸지 않음 — 기존 row가 비어있는 상태)
	if err := db.Exec(`
		ALTER TABLE tasks ADD COLUMN IF NOT EXISTS start_at timestamptz;
		ALTER TABLE tasks ADD COLUMN IF NOT EXISTS end_at   timestamptz;
	`).Error; err != nil {
		return err
	}

	// 2) 기존 row 백필: date 하루를 00:00~23:59로 채움
	result := db.Exec(`
		UPDATE tasks
		SET start_at = date_trunc('day', date),
		    end_at   = date_trunc('day', date) + interval '23 hours 59 minutes 59 seconds'
		WHERE start_at IS NULL;
	`)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected > 0 {
		slog.Info("backfilled task schedule columns", "rows", result.RowsAffected)
	}

	// 3) 백필이 끝났으니 NOT NULL 제약 확정 (모델 태그와 실제 DB 제약을 일치시킴)
	if err := db.Exec(`
		ALTER TABLE tasks ALTER COLUMN start_at SET NOT NULL;
		ALTER TABLE tasks ALTER COLUMN end_at   SET NOT NULL;
	`).Error; err != nil {
		return err
	}

	return nil
}
