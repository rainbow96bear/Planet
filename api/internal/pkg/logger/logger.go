package logger

import (
	"log/slog"
	"os"
)

// New는 slog.Logger를 생성한다.
// 포맷은 JSON 고정 — Cloud Run(Cloud Logging)이 JSON 로그를 자동 파싱해서
// severity/필드 검색을 지원하기 때문. 로컬에서 보기 불편하면 `| jq` 등으로 확인.
func New() *slog.Logger {
	level := slog.LevelInfo
	if os.Getenv("LOG_LEVEL") == "debug" {
		level = slog.LevelDebug
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Cloud Logging은 "level"이 아니라 "severity" 필드로 로그 레벨을 인식한다.
			// 이 매핑이 없으면 Cloud Run 콘솔에서 ERROR/WARNING 필터링이 동작하지 않는다.
			if a.Key == slog.LevelKey {
				a.Key = "severity"
				lvl, _ := a.Value.Any().(slog.Level)
				switch {
				case lvl >= slog.LevelError:
					a.Value = slog.StringValue("ERROR")
				case lvl >= slog.LevelWarn:
					a.Value = slog.StringValue("WARNING")
				case lvl >= slog.LevelInfo:
					a.Value = slog.StringValue("INFO")
				default:
					a.Value = slog.StringValue("DEBUG")
				}
			}
			return a
		},
	})

	return slog.New(handler)
}
