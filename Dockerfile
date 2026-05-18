# Stage 1: Build stage
FROM golang:1.25.0-alpine AS builder

# Install tzdata and ca-certificates for secure requests and correct timezones
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# Copy dependency files and download
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the statically linked Go biner for maximum security and portability
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/go-beresin-api cmd/api/main.go

# Stage 2: Runtime stage
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Copy compiled binary from builder stage
COPY --from=builder /app/go-beresin-api .

# Copy runtime resource folders
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/docs ./docs

# Create uploads folder for file upload support
RUN mkdir -p uploads

# Expose Fiber port
EXPOSE 8080

# Command to run
ENTRYPOINT ["./go-beresin-api"]
