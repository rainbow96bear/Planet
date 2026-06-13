# ================================
# Stage 1: Build
# ================================
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

# go.mod의 버전 선언을 1.24로 변경 (로컬 파일은 영향 없음)
RUN sed -i 's/^go 1\.25.*/go 1.24/' go.mod

RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -o server ./cmd/server/main.go

# ================================
# Stage 2: Runtime
# ================================
FROM alpine:3.21

RUN addgroup -S appgroup && adduser -S appuser -G appgroup
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app
COPY --from=builder /app/server .

RUN chown -R appuser:appgroup /app
USER appuser

EXPOSE 8080

ENTRYPOINT ["./server"]