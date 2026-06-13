# ================================
# Stage 1: Build
# ================================
FROM golang:1.24-alpine AS builder

# 빌드에 필요한 패키지 설치 (CGO 사용 시 gcc 추가)
RUN apk add --no-cache git

WORKDIR /app

# 의존성 캐시 레이어 분리 (소스 변경 시 재다운로드 방지)
COPY go.mod go.sum ./
RUN go mod download

# 소스 복사 및 빌드
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -o server ./cmd/main.go

# ================================
# Stage 2: Runtime
# ================================
FROM alpine:3.21 AS runner

# 보안: 루트 실행 방지용 유저 생성
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# HTTPS 요청, 타임존 데이터
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# 빌드 결과물만 복사
COPY --from=builder /app/server .

# 필요하다면 마이그레이션 파일도 복사
# COPY --from=builder /app/migrations ./migrations

# 소유권 변경 후 비루트 유저로 실행
RUN chown -R appuser:appgroup /app
USER appuser

EXPOSE 8080

ENTRYPOINT ["./server"]