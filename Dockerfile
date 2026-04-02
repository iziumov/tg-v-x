FROM golang:1.26-alpine AS builder

WORKDIR /app

# Install git (needed for go mod)
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

# Runtime stage
FROM alpine:latest

# Install ffmpeg, curl and ca-certificates (needed for yt-dlp and HTTPS requests)
RUN apk add --no-cache ffmpeg curl ca-certificates

# Download yt-dlp binary
RUN curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -o /usr/local/bin/yt-dlp && \
    chmod a+rx /usr/local/bin/yt-dlp

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/main .

# Copy .env file (optional, you can also use docker-compose env)
COPY .env .

CMD ["./main"]
