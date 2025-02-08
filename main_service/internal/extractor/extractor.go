package extractor

import (
    "bytes"
    "encoding/json"
    "net/http"
)

// ExtractResponse represents the response from the Audio Extractor service
type ExtractResponse struct {
    Status    string `json:"status"`
    AudioPath string `json:"audio_path,omitempty"`
    Error     string `json:"error,omitempty"`
}

// Client represents a client implementation for the Audio Extractor
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

// Extract sends a request to the Audio Extractor service to extract audio from video
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