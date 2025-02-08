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

	// Валидация входных данных
	if req.AudioPath == "" || req.TextPath == "" {
		http.Error(w, `{"error":"Missing audio_path or text_path"}`, http.StatusBadRequest)
		return
	}

	// Команда для транскрибирования с помощью whisper
	// Добавлено: указание языка ru
	cmd := exec.Command("whisper", req.AudioPath, "--model", "base", "--output_format", "txt", "--output_dir", filepath.Dir(req.TextPath), "--language", "ru")
	if err := cmd.Run(); err != nil {
		resp := TranscribeResponse{
			Status: "error",
			Error:  err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
		return
	}

	// Читаем транскрибированный текст
	text, err := os.ReadFile(req.TextPath)
	if err != nil {
		resp := TranscribeResponse{
			Status: "error",
			Error:  err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp := TranscribeResponse{
		Status: "success",
		Text:   string(text),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	http.HandleFunc("/transcribe", transcribeHandler)
	log.Println("Transcriber Service is running on port 5002...")
	if err := http.ListenAndServe(":5002", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
