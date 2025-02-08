package main

import (
	"log"
	"net/http"

	"transcriber_project/main_service/internal/downloader"
	"transcriber_project/main_service/internal/extractor"
	"transcriber_project/main_service/internal/transcriber"
	"transcriber_project/main_service/internal/utils"
)

const (
	resultsDir    = "/app/results"
	videosDir     = resultsDir + "/videos"
	outputCSVPath = resultsDir + "/output.csv"
	downloadLog   = resultsDir + "/downloaded.log"
)

func main() {
	log.Println("Starting Main Service...")

	// Initialize service clients
	downloaderClient := downloader.New("http://downloader_container:5011", videosDir, downloadLog)
	extractorClient := extractor.New("http://audio_extractor_container:5001")
	transcriberClient := transcriber.New("http://transcriber_container:5002")

	// Set up HTTP server
	http.HandleFunc("/process-videos", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Processing videos...")

		// Process video files
		utils.ProcessVideos(extractorClient, transcriberClient, downloaderClient, videosDir, outputCSVPath)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Videos processed successfully"))
	})

	// Start HTTP server on port 8080
	log.Fatal(http.ListenAndServe(":8080", nil))
}
