package extractor

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// ExtractResponse представляет ответ от сервиса Audio Extractor
type ExtractResponse struct {
	Status    string `json:"status"`
	AudioPath string `json:"audio_path,omitempty"`
	Error     string `json:"error,omitempty"`
}

// Client представляет собой реализацию клиента для Audio Extractor
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewClient создаёт новый экземпляр Client
func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{},
	}
}

// Extract выполняет запрос к сервису Audio Extractor для извлечения аудио из видео
func (c *Client) Extract(videoPath, audioPath string) (*ExtractResponse, error) {
	requestBody, err := json.Marshal(map[string]string{
		"video_path": videoPath,
		"audio_path": audioPath,
	})
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Post(c.BaseURL+"/extract", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var extractResp ExtractResponse
	if err := json.NewDecoder(resp.Body).Decode(&extractResp); err != nil {
		return nil, err
	}

	return &extractResp, nil
}
