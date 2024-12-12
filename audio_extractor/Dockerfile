# audio_extractor/Dockerfile

# Этап сборки
FROM golang:1.20 AS builder

WORKDIR /app

# Копируем исходный код
COPY . .

# Собираем бинарник
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o audio_extractor main.go

# Финальный образ
FROM alpine:latest

WORKDIR /app

# Устанавливаем ffmpeg
RUN apk add --no-cache ffmpeg

# Копируем бинарник из билдера
COPY --from=builder /app/audio_extractor .

# Устанавливаем права на выполнение
RUN chmod +x ./audio_extractor

# Определяем порт
EXPOSE 5001

# Запускаем сервис
CMD ["./audio_extractor"]
