package utils

import (
    "encoding/csv"
    "log"
    "os"
    "path/filepath"
    "strings"

    "transcriber_project/main_service/internal/downloader"
    "transcriber_project/main_service/internal/extractor"
    "transcriber_project/main_service/internal/transcriber"
)

// ProcessVideos processes video files, extracts audio, transcribes, and writes results to a CSV
func ProcessVideos(
    extractorClient *extractor.Client,
    transcriberClient *transcriber.Client,
    downloaderClient *downloader.Client,
    videosDir, outputCSVPath string,
) {
    err := downloaderClient.DownloadVideosFromFile(videosDir + "/download.txt")
    if err != nil {
        log.Printf("Failed to download video files from file: %v", err)
    }

    // Create the results directory if it doesn't exist
    resultsDir := filepath.Join(videosDir, "results")
    if err := os.MkdirAll(resultsDir, 0755); err != nil {
        log.Fatalf("Failed to create results directory: %v", err)
    }

    // Open (or create) the CSV file for writing results
    file, err := os.Create(outputCSVPath)
    if err != nil {
        log.Fatalf("Failed to create CSV file: %v", err)
    }
    defer file.Close()

    writer := csv.NewWriter(file)
    defer writer.Flush()

    // Write the header
    if err := writer.Write([]string{"Video File", "Transcribed Text"}); err != nil {
        log.Fatalf("Failed to write header to CSV: %v", err)
    }

    // Get the list of video files
    videoFiles, err := os.ReadDir(videosDir)
    if err != nil {
        log.Fatalf("Failed to read videos directory: %v", err)
    }

    for _, video := range videoFiles {
        if video.IsDir() {
            continue
        }

        videoPath := filepath.Join(videosDir, video.Name())
        audioPath := filepath.Join(videosDir, strings.TrimSuffix(video.Name(), filepath.Ext(video.Name()))+".wav")
        textPath := filepath.Join(resultsDir, strings.TrimSuffix(video.Name(), filepath.Ext(video.Name()))+".txt")

        // Assume methods are called Extract and Transcribe
        extractResp, err := extractorClient.Extract(videoPath, audioPath)
        if err != nil || extractResp.Status != "success" {
            log.Printf("Error extracting audio from %s: %v", video.Name(), extractResp.Error)
            continue
        }

        transcribeResp, err := transcriberClient.Transcribe(audioPath, textPath)
        if err != nil || transcribeResp.Status != "success" {
            log.Printf("Error transcribing %s: %v", video.Name(), transcribeResp.Error)
            continue
        }

        // Read the transcribed text
        transcribedTextBytes, err := os.ReadFile(textPath)
        if err != nil {
            log.Printf("Failed to read file %s: %v", textPath, err)
            continue
        }
        transcribedText := string(transcribedTextBytes)

        // Process the text:
        // 1. Replace newlines with spaces
        // 2. Remove quotes
        // 3. Replace commas with semicolons (to prevent CSV conflicts)
        processedText := strings.ReplaceAll(transcribedText, "\r\n", " ")
        processedText = strings.ReplaceAll(processedText, "\n", " ")
        processedText = strings.ReplaceAll(processedText, "\"", "")
        processedText = strings.ReplaceAll(processedText, ",", ";")

        // Write the result to CSV
        record := []string{video.Name(), processedText}
        if err := writer.Write(record); err != nil {
            log.Printf("Error writing record to CSV for file %s: %v", video.Name(), err)
            continue
        }
    }

    // Check for errors during writing
    if err := writer.Error(); err != nil {
        log.Fatalf("Error writing to CSV: %v", err)
    }
}