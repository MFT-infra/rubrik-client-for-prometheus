# --- Stage 1: Build the Go binary ---
FROM golang:1.22-alpine AS builder

# Installs CA certificates and git for fetching dependencies
RUN apk add --no-cache git ca-certificates

WORKDIR /src/golang

# Copy go.mod and go.sum files to leverage Docker cache
# This allows us to avoid re-downloading dependencies unless they change
COPY go.mod go.sum ./
RUN go mod download
RUN go mod tidy
RUN go mod verify

# Copy source code
COPY . .

# Build binary
RUN go build -o /app/rubrik_exporter .

# --- Stage 2: Minimal runtime image ---
FROM alpine:3.20

# Install CA certificates
# This is necessary for HTTPS requests made by the exporter
RUN apk add --no-cache ca-certificates

# Optional: add a non-root user
RUN addgroup -S rubrik && adduser -S rubrik -G rubrik

# Copy binary
COPY --from=builder /app/rubrik_exporter /app/rubrik_exporter

# Set working dir
WORKDIR /app

# Use non-root user
USER rubrik

# Expose Prometheus metrics port
EXPOSE 8080

# Entry point
ENTRYPOINT ["/app/rubrik_exporter"]
