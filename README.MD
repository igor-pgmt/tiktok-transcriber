# Tiktok Transcriber

## Overview

This project is designed for automatic bulk downloading and transcribing of videos from TikTok. It provides a streamlined process to handle large volumes of video content, extracting audio and converting it into text format for further analysis or use.

## Quick Start

1. **Prepare Input**: 
   Fill the `results/download.txt` file with your TikTok video links, one per line.

2. **Start the Service**: 
  Run the following command in your terminal:
```
make run
```

3. **Process Videos**: 
Send a POST request to initiate the processing:
```
curl -X POST http://localhost:8080/process-videos
```

## Requirements

- Docker
- Make

## Project Structure

- `config.env`: Config file for servers ports
- `results/download.txt`: Input file for TikTok video links
- `results/videos/`: Directory where downloaded videos are stored
- `results/output.csv`: Output file containing transcription results

## Additional Information

Currently, the project is single-threaded and operates strictly sequentially.

The project uses kalandar5862/video-downloader as the video file downloader.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
