package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("TRANSCRIBER_PORT")
	if port == "" {
		port = "5002" // Default port if not set
	}

	http.HandleFunc("/transcribe", transcribeHandler)

	log.Printf("Main Service is running on port %s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
