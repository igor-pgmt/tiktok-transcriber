package downloader

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Client represents a client implementation for downloading videos
type Client struct {
	baseURL       string
	downloadDir   string
	downloadedLog string
}

// New creates a new instance of Client
func New(baseURL, downloadDir, downloadedLog string) *Client {
	return &Client{
		baseURL:       baseURL,
		downloadDir:   downloadDir,
		downloadedLog: downloadedLog,
	}
}

// DownloadVideo downloads a video from the specified URL
func (c *Client) DownloadVideo(videoURL string) error {
	if c.isAlreadyDownloaded(videoURL) {
		return nil
	}

	requestData := fmt.Sprintf("url=%s", videoURL)

	// Create a new POST request
	req, err := http.NewRequest("POST", c.baseURL+"/download", bytes.NewBufferString(requestData))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	setHeaders(req)

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error executing request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("failed to download video")
	}

	fileName := filepath.Base(videoURL) + getFileExtension(resp.Header.Get("Content-Type"))
	filePath := filepath.Join(c.downloadDir, fileName)
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("error writing file: %v", err)
	}

	return c.logDownloaded(videoURL)
}

// isAlreadyDownloaded checks if the video has already been downloaded
func (c *Client) isAlreadyDownloaded(videoURL string) bool {
	file, err := os.Open(c.downloadedLog)
	if err != nil {
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if scanner.Text() == videoURL {
			return true
		}
	}

	return false
}

// logDownloaded logs the URL of the downloaded video
func (c *Client) logDownloaded(videoURL string) error {
	file, err := os.OpenFile(c.downloadedLog, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(videoURL + "\n")
	if err != nil {
		return fmt.Errorf("failed to write to log file: %v", err)
	}

	return nil
}

// DownloadVideosFromFile downloads videos from a file containing a list of URLs
func (c *Client) DownloadVideosFromFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	var totalLines int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.TrimSpace(scanner.Text()) != "" {
			totalLines++
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	// Second pass to process lines
	file.Seek(0, io.SeekStart) // Reset file pointer to the beginning
	scanner = bufio.NewScanner(file)
	var currentLine int

	for scanner.Scan() {
		videoURL := strings.TrimSpace(scanner.Text())
		if videoURL == "" {
			continue
		}

		err := c.DownloadVideo(videoURL)
		if err != nil {
			fmt.Printf("Error downloading video %s: %v\n", videoURL, err)
		}

		currentLine++
		progress := float64(currentLine) / float64(totalLines) * 100
		fmt.Printf("Progress %d of %d: %.2f%%\n", currentLine, totalLines, progress)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	return nil
}

func getFileExtension(contentType string) string {
	switch contentType {
	case "video/mp4":
		return ".mp4"
	case "video/x-matroska":
		return ".mkv"
	case "video/webm":
		return ".webm"
	case "video/ogg":
		return ".ogv"
	case "video/quicktime":
		return ".mov"
	case "video/avi":
		return ".avi"
	case "video/mpeg":
		return ".mpeg"
	default:
		return ""
	}
}

func setHeaders(req *http.Request) {
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,ka;q=0.8")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Length", "79")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Host", "localhost:5011")
	req.Header.Set("Origin", "http://localhost:5011")
	req.Header.Set("Referer", "http://localhost:5011/")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36")
	req.Header.Set("sec-ch-ua", `"Google Chrome";v="131", "Chromium";v="131", "Not_A Brand";v="24"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "macOS")
}
