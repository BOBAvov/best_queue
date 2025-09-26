# Используем официальный образ Golang для сборки
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Копируем go.mod и go.sum для установки зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем бинарник
RUN go build -o http_sso ./cmd/main.go

# Финальный минимальный образ
FROM alpine:latest

WORKDIR /app

# Копируем бинарник из стадии сборки
COPY --from=builder /app/http_sso .

# Открываем порт
EXPOSE 8080

# Запускаем приложение
CMD ["./http_sso"]