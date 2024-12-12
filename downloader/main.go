package main

import (
	"fmt"
	"log"

	"github.com/pawanpaudel93/go-tiktok-downloader/tiktok"
)

func main() {
	videoURL := "https://www.tiktok.com/@iakoo83/video/7422377711680736520"

	video := tiktok.Video{URL: videoURL, BaseDIR: ".", Proxy: ""}

	videoPath, err := video.Download()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Видео успешно скачано по пути: %s\n", videoPath)

	// Извлечение метаданных
	metadata, err := video.GetInfo()
	if err != nil {
		log.Fatalf("Ошибка при получении метаданных: %v", err)
	}

	fmt.Printf("Метаданные видео: %+v\n", metadata)
}
