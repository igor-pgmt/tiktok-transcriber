package main

import (
	"log"
	"net/http"

	"transcriber_project/main_service/internal/downloader"
	"transcriber_project/main_service/internal/extractor"
	"transcriber_project/main_service/internal/transcriber"
	"transcriber_project/main_service/internal/utils"
)

func NewHandler(downloaderClient *downloader.Client, extractorClient *extractor.Client, transcriberClient *transcriber.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Processing videos...")

		// Process video files
		utils.ProcessVideos(extractorClient, transcriberClient, downloaderClient, videosDir, fileToDownload, outputCSVPath)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Videos processed successfully"))
	}
}
