package main

import (
	"fmt"
	"log"
	"net/http"
	"transcriber_project/main_service/internal/downloader"
	"transcriber_project/main_service/internal/extractor"
	"transcriber_project/main_service/internal/transcriber"
)

const (
	resultsDir    = "/app/results"
	videosDir     = resultsDir + "/videos"
	outputCSVPath = resultsDir + "/output.csv"
	downloadLog   = resultsDir + "/downloaded.log"
)

func main() {
	cfg := newConfig()

	// Initialize service clients
	downloaderClient := downloader.New(fmt.Sprintf("http://downloader_container:%s", cfg.DownloaderPort), videosDir, downloadLog)
	extractorClient := extractor.New(fmt.Sprintf("http://audio_extractor_container:%s", cfg.AudioExtractorPort))
	transcriberClient := transcriber.New(fmt.Sprintf("http://transcriber_container:%s", cfg.TranscriberPort))

	mainHandler := NewHandler(downloaderClient, extractorClient, transcriberClient)

	// Set up HTTP server
	http.HandleFunc("/process-videos", mainHandler)

	log.Printf("Main Service is running on port %s...", cfg.MainServicePort)
	if err := http.ListenAndServe(":"+cfg.MainServicePort, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
