# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app
RUN apk add --no-cache git

# Copy only the needed folders
COPY common-service/ ./common-service/
COPY db/ ./db/

# Set the working directory to the service
WORKDIR /app/common-service

# Build the service
RUN go build -o common-service .

# Final stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/common-service/common-service .
COPY --from=builder /app/common-service/.env .env
EXPOSE 8080
CMD ["./common-service"]

