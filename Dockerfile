# Use Go base image
FROM golang:1.20-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go modules and dependencies
COPY go.mod go.sum ./
RUN go mod tidy

# Copy the rest of the application
COPY . .

# Pass environment variables as build arguments
ARG DATABASE_URL
ARG PORT
ARG JWTKEY
ARG REDISHOST
ARG REDIS_PASSWORD

# Build the application
RUN go build -o task-mgt-crud-app ./cmd

# Create a lightweight final image
FROM alpine:latest
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/task-mgt-crud-app .

# Set environment variables
ENV DATABASE_URL=$DATABASE_URL
ENV PORT=$PORT
ENV JWTKEY=$JWTKEY
ENV REDISHOST=$REDISHOST
ENV REDIS_PASSWORD=$REDIS_PASSWORD

# Expose the application port
EXPOSE 8080

# Start the application
CMD ["./task-mgt-crud-app"]



# #build stage - stage-1
# FROM golang:1.23.1-alpine AS builder 

# WORKDIR /app

# #copy source code
# COPY . /app

# #build the application binary code output file
# RUN  go build -o task-mgt-crud-app ./cmd

# #production stage -stage -2
# FROM alpine:latest

# WORKDIR /app

# #Copy the built binary from the builder stage
# COPY --from=builder /app/task-mgt-crud-app .

# COPY .env /app/ 
# # Expose the application port
# EXPOSE 8080

# # Command to run the application
# CMD ["./task-mgt-crud-app"]