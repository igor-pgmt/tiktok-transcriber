package main

import (
	"encoding/json"
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

	// Input validation
	if req.VideoPath == "" || req.AudioPath == "" {
		http.Error(w, `{"error":"Missing video_path or audio_path"}`, http.StatusBadRequest)
		return
	}

	// Command to extract audio using ffmpeg
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
