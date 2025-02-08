# audio_extractor/Dockerfile

# Build stage
FROM golang:1.23 AS builder

WORKDIR /app

# Copy the source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o audio_extractor *.go

# Final image
FROM alpine:latest

WORKDIR /app

# Install ffmpeg
RUN apk add --no-cache ffmpeg

# Copy the binary from the builder
COPY --from=builder /app/audio_extractor .

# Set execution permissions
RUN chmod +x ./audio_extractor

# Run the service
CMD ["./audio_extractor"]