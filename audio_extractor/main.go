package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("AUDIO_EXTRACTOR_PORT")
	if port == "" {
		port = "5001" // Default port if not set
	}

	http.HandleFunc("/extract", extractHandler)

	log.Printf("Audio Extractor Service is running on port %s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
