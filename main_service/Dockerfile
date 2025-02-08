# main_service/Dockerfile

# Use the official Go image for building
FROM golang:1.23 AS builder

WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/*.go

# Final image
FROM alpine:latest

# Changed: set the working directory to /app/
WORKDIR /app/

# Install necessary libraries
RUN apk add --no-cache ca-certificates

# Copy the binary from the builder
COPY --from=builder /app/main .

# Run the binary
CMD ["./main"]