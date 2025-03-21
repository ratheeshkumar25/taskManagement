# Stage 1: Build the Go application
FROM golang:1.23.1-alpine AS builder 

# Set working directory
WORKDIR /app

# Copy go modules and dependencies first (for caching)
COPY go.mod go.sum ./
RUN go mod tidy

# Copy the application source code
COPY . .

# Build the application binary
RUN go build -o task-mgt-crud-app ./cmd

# Stage 2: Create a lightweight production image
FROM alpine:3.18

# Install required packages
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/task-mgt-crud-app .

# Copy environment file for configuration
COPY .env /app/

# List files for debugging
RUN ls -la /app

# Expose the application port
EXPOSE 8080

# Command to run the application
CMD ["./task-mgt-crud-app"]
