# 로깅 가이드

Planet 백엔드(Go)에서 로그를 남길 때 따르는 규칙: **패키지 선정**, **레이어 원칙**, **기록 양식**, **레벨 사용 기준**.

## 0. 로그 패키지

- **사용 패키지**: `log/slog` (Go 1.21+ 표준 라이브러리)
- **선정 이유**
  - 표준 라이브러리라 외부 의존성 추가 없이 바로 사용 가능
  - 구조화된 key-value 로깅(`slog.Info("msg", "key", value)`)이 언어 차원에서 지원되어 별도 래퍼가 필요 없음
  - `slog.NewJSONHandler`로 JSON 출력을 기본 제공 — Cloud Run(Cloud Logging)이 JSON 로그를 자동 파싱해서 필드 검색/필터링 지원
  - `slog.SetDefault`로 전역 로거를 한 번만 설정하면 패키지 함수(`slog.Info` 등)로 어디서든 동일한 설정을 사용
- **구성**: `internal/pkg/logger.New()`가 JSON 핸들러를 생성한다.
  - `LOG_LEVEL` 환경변수로 레벨 조정 (`debug` 지정 시 Debug까지, 기본은 Info 이상)
  - Cloud Logging이 인식하는 `severity` 필드로 slog의 `level`을 매핑 (`ERROR`/`WARNING`/`INFO`/`DEBUG`) — 이 매핑이 없으면 Cloud Run 콘솔의 심각도 필터가 동작하지 않음
  - `main.go`에서 앱 시작 시 `slog.SetDefault(logger.New())`로 한 번만 초기화

## 1. 레이어 원칙 — 어디에 로그를 남기나

| 레이어 | 로그 | 이유 |
|---|---|---|
| repository | ❌ | DB 에러를 그대로 위로 반환. 여기서 찍으면 service에서 같은 에러가 또 찍혀 중복됨 |
| **service** | ✅ 유일한 위치 | 비즈니스 판단(성공/실패, 권한, 데이터 정합성)이 일어나는 곳 |
| handler | ❌ | `slogRequestLogger` 미들웨어가 요청 단위(method/path/status/latency)로 이미 커버 |

**동일 오류 중복 기록 방지 규칙**
- 하나의 에러는 그 에러를 "처리"하는 계층에서 딱 한 번만 로그한다. 처리한다는 건 로그를 남기고 그대로 반환하거나, 사용자용 메시지로 바꿔서 반환하는 것을 말한다.
- repository에서 발생한 에러가 service까지 그대로 올라오면 service에서 로그 + 반환. 이후 handler는 그 에러를 HTTP 상태코드로 변환만 하고 다시 로그하지 않는다.
- 새 코드에서 로그를 넣고 싶은데 이미 하위 계층에서 같은 에러를 반환받아 로그 없이 그대로 위로 넘기고 있다면, 그 지점(가장 처음 에러를 받는 service 메서드)에서만 로그하고 그 이후 계층에서는 추가하지 않는다.

## 2. 기록 양식

```go
slog.레벨("무엇을 하려다 성공/실패했는지", "관련_id", 값, ..., "error", err)
```

**규칙**
- 메시지는 영어 소문자, 동사 원형으로 작성한다. 예: `"failed to create task"`, `"login success"`, `"unauthorized delete attempt"`
- 필드 순서: 메시지 → 관련 엔티티 ID들 → (실패 시) 마지막에 `"error", err`
- 필드 이름은 프로젝트 전체에서 통일해서 재사용한다. 새로 만들지 말고 아래 표에서 고른다.

| 필드명 | 용도 |
|---|---|
| `error` | 에러 값 (항상 이 이름 사용, `err` 금지) |
| `user_id` | 요청/작업 주체의 사용자 ID |
| `requester_id` | 조회를 요청한 사람 (대상과 다를 때) |
| `receiver_id` / `actor_id` | 알림 등에서 받는 사람 / 행동한 사람 |
| `task_id`, `notification_id` 등 | 관련 리소스 ID |

**예시**

```go
// 실패
if err := s.taskRepo.CreateTask(tx, task); err != nil {
	slog.Error("failed to create task", "user_id", req.UserID, "error", err)
	return nil, err
}

// 성공 (error 필드 없음)
slog.Info("task created", "task_id", task.ID, "user_id", req.UserID)

// 보안/권한 관련
slog.Warn("unauthorized delete attempt",
	"task_id", req.ID,
	"owner_id", task.UserID,
	"requester_id", req.UserID,
)
```

**절대 로그에 남기지 않는 것**: 비밀번호(평문/해시 모두), Access/Refresh Token 전체 값, OAuth Secret, 쿠키/세션 값, 그 외 개인정보.

## 3. 레벨별 사용 기준

| 레벨 | 언제 쓰나 | 예시 |
|---|---|---|
| **Error** | 시스템/인프라 실패. "우리 잘못" | DB 에러, 트랜잭션 실패, 외부 스토리지 I/O 실패, 토큰 생성 실패 |
| **Warn** | 복구는 되지만 주목할 상황 | 로그인 실패, 권한 위반, best-effort 작업 실패(실패해도 전체 흐름은 유지), 원인이 뭉뚱그려진 에러 |
| **Info** | 정상적이고 의미 있는 이벤트 (성공 시) | 회원가입, 로그인 성공, 생성/삭제/수정 성공, 서버 시작/종료 |
| **Debug** | 로컬 개발용 상세 정보 | 필요해질 때만 추가, 운영에서는 기본적으로 안 보임 |

**주의**: 고빈도 조회(피드/검색/알림 폴링)나 자주 반복되는 쓰기 액션(좋아요, 팔로우, 완료 토글)은 **성공 로그를 남기지 않는다** — Error/Warn만. 매번 찍으면 로그가 순식간에 노이즈가 된다. 가입·로그인·삭제처럼 드물고 감사 가치가 있는 액션만 성공까지 Info로 남긴다.

## 4. JSON 로그 출력 형식

Cloud Run에 찍히는 로그 한 줄은 대략 다음 필드를 가진다 (slog JSON 핸들러 기본 + `severity` 매핑):

```json
{
  "time": "2026-09-14T12:34:56.789Z",
  "severity": "ERROR",
  "msg": "failed to create task",
  "user_id": "u_123",
  "error": "duplicate key value violates unique constraint"
}
```

## 5. 향후 개선 (미구현)

- **요청 상관관계(request_id)**: 현재 `slogRequestLogger`가 찍는 요청 단위 로그와 service에서 찍는 비즈니스 로그를 하나의 요청으로 묶어줄 공통 식별자가 없음. Gin 미들웨어에서 요청마다 ID를 발급해 context로 전달하고, 모든 로그에 `request_id` 필드를 포함시키는 구조를 추후 도입 예정.