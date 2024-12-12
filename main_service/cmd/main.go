package main

import (
	"log"
	"transcriber_project/main_service/internal/extractor"
	"transcriber_project/main_service/internal/transcriber"
	"transcriber_project/main_service/internal/utils"
)

func main() {
	log.Println("Starting Main Service...")

	// Инициализация клиентов сервисов
	extractorClient := extractor.NewClient("http://audio_extractor_container:5001")
	transcriberClient := transcriber.NewClient("http://transcriber_container:5002")

	// Обработка видеофайлов
	// Изменено: сохраняем output.csv в /app/results/
	utils.ProcessVideos("/app/videos", extractorClient, transcriberClient, "/app/results/output.csv")
}
