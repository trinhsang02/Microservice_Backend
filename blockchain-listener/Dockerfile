# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app
RUN apk add --no-cache git

# Copy only the needed folders
COPY blockchain-listener/ ./blockchain-listener/
COPY db/ ./db/
COPY mq/ ./mq/

# Set the working directory to the service
WORKDIR /app/blockchain-listener

# Build the service
RUN go build -o blockchain-listener .

# Final stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/blockchain-listener/blockchain-listener .
COPY --from=builder /app/blockchain-listener/.env .env
EXPOSE 8081
EXPOSE 50051
CMD ["./blockchain-listener"]