#build stage - stage-1
FROM golang:1.23.1-alpine AS builder 

WORKDIR /app

#copy source code
COPY . /app

#build the application binary code output file
RUN  go build -o task-mgt-crud-app ./cmd

#production stage -stage -2
FROM alpine:latest

WORKDIR /app

#Copy the built binary from the builder stage
COPY --from=builder /app/task-mgt-crud-app .

COPY .env /app/ 
# Expose the application port
EXPOSE 8080

# Command to run the application
CMD ["./task-mgt-crud-app"]