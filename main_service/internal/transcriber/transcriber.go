package transcriber

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// TranscribeResponse представляет ответ от сервиса Transcriber
type TranscribeResponse struct {
	Status string `json:"status"`
	Text   string `json:"text,omitempty"`
	Error  string `json:"error,omitempty"`
}

// Client представляет собой реализацию клиента для Transcriber
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

// Transcribe выполняет запрос к сервису Transcriber для транскрипции аудио
func (c *Client) Transcribe(audioPath, textPath string) (*TranscribeResponse, error) {
	requestBody, err := json.Marshal(map[string]string{
		"audio_path": audioPath,
		"text_path":  textPath,
	})
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Post(c.BaseURL+"/transcribe", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var transcribeResp TranscribeResponse
	if err := json.NewDecoder(resp.Body).Decode(&transcribeResp); err != nil {
		return nil, err
	}

	return &transcribeResp, nil
}
