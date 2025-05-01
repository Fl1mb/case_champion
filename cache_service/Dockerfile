# Используем образ Go на основе Alpine
FROM golang:1.23.2-alpine AS builder

# Устанавливаем зависимости через apk
RUN apk add --no-cache \
    protoc \
    protobuf-dev \
    git \
    gcc \
    musl-dev

# Устанавливаем Go-плагины для protoc
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest \
    && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /app/bin/myapp ./cmd/main.go

# Финальный образ (Alpine)
FROM alpine:latest
COPY --from=builder /app/bin/myapp /usr/local/bin/myapp
ENTRYPOINT ["/usr/local/bin/myapp"]