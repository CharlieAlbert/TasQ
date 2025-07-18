# Stage 1: Build
FROM golang:1.23 AS builder

WORKDIR /app

# Copy dependency files first (for caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the code
COPY . .

# Build the binary with proper flags for Alpine
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/job_dispatcher ./cmd/main.go

# Stage 2: Runtime
FROM alpine:3.20

WORKDIR /app

# Install CA certificates (needed for TLS) and libc6-compat (for Go binaries)
RUN apk update && apk add --no-cache ca-certificates libc6-compat

# Copy the compiled binary from builder
COPY --from=builder /app/job_dispatcher ./job_dispatcher

# Make the binary executable
RUN chmod +x ./job_dispatcher

EXPOSE 8080

# Run the binary
CMD ["./job_dispatcher"]