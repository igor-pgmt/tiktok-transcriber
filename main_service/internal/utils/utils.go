package utils

import (
	"encoding/csv"
	"log"
	"os"
	"path/filepath"
	"strings"

	"transcriber_project/main_service/internal/extractor"
	"transcriber_project/main_service/internal/transcriber"
)

// ProcessVideos обрабатывает видеофайлы, извлекает аудио, транскрибирует и записывает результаты в CSV
func ProcessVideos(videosDir string, extractorClient *extractor.Client, transcriberClient *transcriber.Client, outputCSVPath string) {
	// Создаём директорию results, если она не существует
	resultsDir := filepath.Join(videosDir, "results")
	if err := os.MkdirAll(resultsDir, 0755); err != nil {
		log.Fatalf("Не удалось создать директорию results: %v", err)
	}

	// Открываем (или создаем) CSV-файл для записи результатов
	file, err := os.Create(outputCSVPath)
	if err != nil {
		log.Fatalf("Не удалось создать файл CSV: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Записываем заголовок
	if err := writer.Write([]string{"Video File", "Transcribed Text"}); err != nil {
		log.Fatalf("Не удалось записать заголовок в CSV: %v", err)
	}

	// Получаем список видеофайлов
	videoFiles, err := os.ReadDir(videosDir)
	if err != nil {
		log.Fatalf("Не удалось прочитать директорию videos: %v", err)
	}

	for _, video := range videoFiles {
		if video.IsDir() {
			continue
		}

		videoPath := filepath.Join(videosDir, video.Name())
		audioPath := filepath.Join(videosDir, strings.TrimSuffix(video.Name(), filepath.Ext(video.Name()))+".wav")
		textPath := filepath.Join(resultsDir, strings.TrimSuffix(video.Name(), filepath.Ext(video.Name()))+".txt")

		// Предположим, методы называются Extract и Transcribe
		extractResp, err := extractorClient.Extract(videoPath, audioPath)
		if err != nil || extractResp.Status != "success" {
			log.Printf("Ошибка при извлечении аудио из %s: %v", video.Name(), extractResp.Error)
			continue
		}

		transcribeResp, err := transcriberClient.Transcribe(audioPath, textPath)
		if err != nil || transcribeResp.Status != "success" {
			log.Printf("Ошибка при транскрибировании %s: %v", video.Name(), transcribeResp.Error)
			continue
		}

		// Чтение транскрибированного текста
		transcribedTextBytes, err := os.ReadFile(textPath)
		if err != nil {
			log.Printf("Не удалось прочитать файл %s: %v", textPath, err)
			continue
		}
		transcribedText := string(transcribedTextBytes)

		// Обработка текста:
		// 1. Замена переносов строк на пробелы
		// 2. Удаление кавычек
		// 3. Замена запятых на точки с запятой (для предотвращения конфликтов в CSV)
		processedText := strings.ReplaceAll(transcribedText, "\r\n", " ")
		processedText = strings.ReplaceAll(processedText, "\n", " ")
		processedText = strings.ReplaceAll(processedText, "\"", "")
		processedText = strings.ReplaceAll(processedText, ",", ";")

		// Записываем результат в CSV
		record := []string{video.Name(), processedText}
		if err := writer.Write(record); err != nil {
			log.Printf("Ошибка при записи записи в CSV для файла %s: %v", video.Name(), err)
			continue
		}
	}

	// Проверяем наличие ошибок при записи
	if err := writer.Error(); err != nil {
		log.Fatalf("Ошибка при записи в CSV: %v", err)
	}
}
