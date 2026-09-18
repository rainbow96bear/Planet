# 🌍 PLANET

**Calendar-based Social Platform**
시간을 기반으로 사람과 사람을 연결하는 캘린더 기반 소셜 플랫폼입니다.
일정을 기록하고, 이를 선택적으로 공개해 다른 사용자와 공유하는 다이어리 + 소셜 네트워크 성격의 서비스입니다.

> People → Time → Relationship
> "함께 보낸 시간"을 관계의 기준으로 삼는다는 철학 위에서, 팔로우 대신 **Orbit / Gravity**라는 자체 용어 체계를 정의해 백엔드 도메인 모델부터 UI 문구까지 일관되게 적용하고 있습니다.

---

## 🧱 기술 스택

| 영역 | 스택 |
|---|---|
| Backend | Go, Gin, GORM, PostgreSQL (Supabase) |
| Frontend | SvelteKit (Svelte 5 runes) |
| Infra / Deploy | Cloud Run (Backend), Cloudflare Pages (Frontend) |
| CI/CD | GitHub Actions (dev / main 브랜치별 빌드 검증) |
| Auth | JWT + SvelteKit 서버 라우트 프록시 패턴 |

---

## 🏗️ 아키텍처

**Backend**
- Repository - Service - Handler 계층형 아키텍처
- Feed 팬아웃(fan-out) 구조 — `Feed` 테이블에 `TaskID` / `ActorID`를 직접 FK로 연결해 조회 성능 확보
- 삭제 정책 이원화 — `Task`는 Soft Delete, `Feed` / `Reaction`은 Hard Delete, 연쇄 삭제는 트랜잭션으로 정합성 보장

**Frontend**
- Svelte 5 runes(`$state` / `$derived` / `$effect`) 기반 반응형 상태 관리
- 낙관적 UI 업데이트(Optimistic UI) + 실패 시 콜백 기반 롤백

**배포**
- Backend는 Cloud Run, Frontend는 Cloudflare Pages로 분리 배포
- `feature` 브랜치 → `dev` PR → `main` PR 순서의 배포 전략, GitHub Actions로 자동 빌드 검증

---

## 🪐 핵심 개념 — Orbit & Gravity

- **Orbit** — 관심 있는 사람의 시간을 내 캘린더로 끌어오는 행위 (팔로우에 대응)
- **Gravity** — 다른 사람을 자연스럽게 끌어당기는 영향력 (팔로워 수에 대응)

팔로워 수가 아니라 "얼마나 많은 사람의 시간에 연결되어 있는가"를 관계의 지표로 삼는다는 컨셉을 반영한 네이밍입니다.

---

## ✨ 진행 현황

### ✅ 구현 완료
- 프로필 이미지 등록
- Follow / Follower 도메인을 Orbit / Gravity로 전면 리네이밍 (백엔드 · DTO · UI)
- Backend 로그 시스템 구축
- Reaction(좋아요 · 응원) API
- CI 파이프라인 (GitHub Actions, dev/main 브랜치 검증)
- 캘린더 일정 Layer — My Schedule / Orbit Schedule 분리, Layer별 독립 ON/OFF, Date Detail(주간 스트립 + 시간순 뷰)
- 일정 수정(Edit) 기능

### 🚧 진행 중
_(현재 진행 중인 작업 없음 — 다음 계획 1번부터 순서대로 착수 예정)_

### 🔮 다음 계획

1. Orbit / Gravity 목록 보기
2. 회원 탈퇴 (Soft Delete 설계)
3. 피드 UI 개선 및 request/response 구조 개선
4. 일정 시작/종료 시간 변경 시에만 Feed 발행되도록 개선
5. 겹치는 시간 찾기
6. 반복 일정
7. 공개 캘린더 링크 공유
8. 주간 리캡

---

## 🐛 기술적으로 다뤄본 문제들

- **Feed 조회 N+1 쿼리 제거** — 좋아요/응원 수를 Feed 항목마다 별도 쿼리로 집계하던 구조를 PostgreSQL의 `FILTER`/`BOOL_OR` 집계 함수로 재설계해, Feed 목록 조회를 단일 쿼리로 통합. Feed 개수에 비례해 늘어나던 쿼리 수를 O(N)에서 O(1)로 줄여 응답 지연을 개선했다.
- **Soft Delete 연쇄 삭제 시 트랜잭션 정합성 확보** — Task 삭제 시 연관된 Feed/Reaction이 별도 쿼리로 순차 삭제되던 구조에서, 중간 실패 시 Task는 삭제됐지만 Feed는 남아 비공개 데이터가 계속 노출되는 버그를 발견. 삭제 로직 전체를 하나의 DB 트랜잭션으로 묶어 부분 삭제로 인한 데이터 정합성 문제를 원천 차단했다.

---

## 🧭 개발 방향

- 최소 기능(MVP)을 빠르게 구현하고, 실제 사용을 통해 문제를 발견하며 지속적으로 리팩토링
- Sprint 단위(1주)로 로드맵을 관리하며 다양한 기술을 실험하고 적용하는 학습 중심 프로젝트

## 🎯 목표

- RESTful API 설계 경험
- SvelteKit 기반 프론트엔드 구조 설계
- Go 서버 아키텍처 설계 및 구현
- 실제 서비스 운영을 고려한 구조 설계 경험 (로깅, 배포, CI/CD 포함)

---

## ▶️ 실행 방법

### Frontend
```bash
npm install       # 패키지 설치
npm run dev       # 빌드 없이 실행
npm run build     # 빌드
```

### Backend
```bash
go run ./cmd/server/main.go   # 빌드 없이 실행
go build ./cmd/server         # 빌드
```

### Docker (PostgreSQL)
```bash
docker compose up -d      # 실행
docker compose down -v    # 삭제
```

---

## 📚 문서

- [로깅 가이드](./docs/LOGGING.md) — 로그 기록 양식, 레벨별 사용 기준
- [Planet Workspace (Notion)](https://planet-sns.notion.site/Planet-Home-ee8201ca907a46d49436121b8a4bfaae) — 기획 · 로드맵 · 브랜드 문서 전체

---

## 📌 프로젝트 성격

완성도 자체보다 **실험, 학습, 개선 과정**에 의미를 둔 개인 프로젝트입니다.
기획부터 백엔드 · 프론트엔드 · 인프라 · 배포까지 1인 개발로 진행하며, Sprint 단위로 기록을 남기고 있습니다.