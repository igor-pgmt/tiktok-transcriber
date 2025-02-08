package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

type TranscribeRequest struct {
	AudioPath string `json:"audio_path"`
	TextPath  string `json:"text_path"`
}

type TranscribeResponse struct {
	Status string `json:"status"`
	Text   string `json:"text,omitempty"`
	Error  string `json:"error,omitempty"`
}

func transcribeHandler(w http.ResponseWriter, r *http.Request) {
	var req TranscribeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid request payload"}`, http.StatusBadRequest)
		return
	}

	// Input validation
	if req.AudioPath == "" || req.TextPath == "" {
		http.Error(w, `{"error":"Missing audio_path or text_path"}`, http.StatusBadRequest)
		return
	}

	// Command for transcription using whisper
	// Added: specifying language as ru
	cmd := exec.Command("whisper", req.AudioPath, "--model", "base", "--output_format", "txt", "--output_dir", filepath.Dir(req.TextPath), "--language", "ru")
	if err := cmd.Run(); err != nil {
		resp := TranscribeResponse{
			Status: "error",
			Error:  err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error encoding JSON response: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}

		return
	}

	// Read the transcribed text
	text, err := os.ReadFile(req.TextPath)
	if err != nil {
		resp := TranscribeResponse{
			Status: "error",
			Error:  err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error encoding JSON response: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}

		return
	}

	resp := TranscribeResponse{
		Status: "success",
		Text:   string(text),
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
