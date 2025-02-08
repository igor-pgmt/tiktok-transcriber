package transcriber

import (
    "bytes"
    "encoding/json"
    "net/http"
)

// TranscribeResponse represents the response from the Transcriber service
type TranscribeResponse struct {
    Status string `json:"status"`
    Text   string `json:"text,omitempty"`
    Error  string `json:"error,omitempty"`
}

// Client represents a client implementation for the Transcriber
type Client struct {
    BaseURL    string
    HTTPClient *http.Client
}

// New creates a new instance of Client
func New(baseURL string) *Client {
    return &Client{
        BaseURL:    baseURL,
        HTTPClient: &http.Client{},
    }
}

// Transcribe sends a request to the Transcriber service to transcribe audio
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