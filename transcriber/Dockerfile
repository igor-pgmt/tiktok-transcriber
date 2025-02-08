# transcriber/Dockerfile

# Build stage
FROM golang:1.23 AS builder

WORKDIR /app

# Copy the source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o transcriber *.go

# Final image
FROM python:3.9-slim

WORKDIR /app

# Install necessary system dependencies
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    ffmpeg \
    build-essential \
    libffi-dev \
    libssl-dev \
    && rm -rf /var/lib/apt/lists/*

# Install pip and upgrade it
RUN pip install --upgrade pip

# Install torch from the official PyTorch repository for CPU
RUN pip install torch --index-url https://download.pytorch.org/whl/cpu

# Install openai-whisper
RUN pip install openai-whisper

# Copy the binary from the builder
COPY --from=builder /app/transcriber .

# Set execution permissions
RUN chmod +x ./transcriber

# Run the service
CMD ["./transcriber"]