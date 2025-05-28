# Build stage
FROM golang:1.24-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app
RUN apk add --no-cache git

# Copy the source code into the container
COPY common-service/ ./common-service/
COPY db/ ./db/
COPY mq/ ./mq/

# Set working directory to the service
WORKDIR /app/common-service

# Build the Go app
RUN go build -o common-service .

# Final stage
FROM alpine:latest
WORKDIR /app

# Create logs directory
RUN mkdir -p /app/logs

COPY --from=builder /app/common-service/common-service .
COPY --from=builder /app/common-service/.env .env
EXPOSE 8080
CMD ["./common-service"]

