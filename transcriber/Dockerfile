# transcriber/Dockerfile

# Этап сборки
FROM golang:1.20 AS builder

WORKDIR /app

# Копируем исходный код
COPY . .

# Собираем бинарник
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o transcriber main.go

# Финальный образ
FROM python:3.9-slim

WORKDIR /app

# Устанавливаем необходимые системные зависимости
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    ffmpeg \
    build-essential \
    libffi-dev \
    libssl-dev \
    && rm -rf /var/lib/apt/lists/*

# Устанавливаем pip и обновляем его
RUN pip install --upgrade pip

# Устанавливаем torch из официального репозитория PyTorch для CPU
RUN pip install torch --index-url https://download.pytorch.org/whl/cpu

# Устанавливаем openai-whisper
RUN pip install openai-whisper

# Копируем бинарник из билдера
COPY --from=builder /app/transcriber .

# Устанавливаем права на выполнение
RUN chmod +x ./transcriber

# Определяем порт
EXPOSE 5002

# Запускаем сервис
CMD ["./transcriber"]
