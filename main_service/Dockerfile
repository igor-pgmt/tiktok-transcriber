# main_service/Dockerfile

# Используем официальный образ Go для сборки
FROM golang:1.20 AS builder

WORKDIR /app

# Копируем go.mod и go.sum
COPY go.mod go.sum ./

# Скачиваем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем бинарник
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/main.go

# Финальный образ
FROM alpine:latest

# Изменено: установка рабочей директории на /app/
WORKDIR /app/

# Устанавливаем необходимые библиотеки
RUN apk add --no-cache ca-certificates

# Копируем бинарник из билдера
COPY --from=builder /app/main .

# Запускаем бинарник
CMD ["./main"]
