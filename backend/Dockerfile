# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install dependencies required for some CGO/sqlite configurations if needed, though this project uses PG
# RUN apk add --no-cache gcc musl-dev

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
# Use CGO_ENABLED=0 to ensure a statically linked binary
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api/main.go

# Final stage
FROM alpine:latest  

WORKDIR /app

# Copy the pre-built binary file from the previous stage
COPY --from=builder /app/main .
# Copy any required static/config files if any (e.g. .env)
# COPY .env .

# Expose port (adjust if your Go app uses a different port, usually 8080)
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
