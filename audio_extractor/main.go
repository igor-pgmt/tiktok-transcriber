package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os/exec"
)

type ExtractRequest struct {
	VideoPath string `json:"video_path"`
	AudioPath string `json:"audio_path"`
}

type ExtractResponse struct {
	Status    string `json:"status"`
	AudioPath string `json:"audio_path,omitempty"`
	Error     string `json:"error,omitempty"`
}

func extractHandler(w http.ResponseWriter, r *http.Request) {
	var req ExtractRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid request payload"}`, http.StatusBadRequest)
		return
	}

	// Валидация входных данных
	if req.VideoPath == "" || req.AudioPath == "" {
		http.Error(w, `{"error":"Missing video_path or audio_path"}`, http.StatusBadRequest)
		return
	}

	// Команда для извлечения аудио с помощью ffmpeg
	cmd := exec.Command("ffmpeg", "-i", req.VideoPath, "-vn", "-acodec", "pcm_s16le", "-ar", "44100", "-ac", "2", req.AudioPath)
	if err := cmd.Run(); err != nil {
		resp := ExtractResponse{
			Status: "error",
			Error:  err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp := ExtractResponse{
		Status:    "success",
		AudioPath: req.AudioPath,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	http.HandleFunc("/extract", extractHandler)
	log.Println("Audio Extractor Service is running on port 5001...")
	if err := http.ListenAndServe(":5001", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
