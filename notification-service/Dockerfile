# Use an official Golang runtime as a parent image
FROM golang:1.24

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the mq module first
COPY notification-service/ ./notification-service/
COPY mq/ ./mq/

WORKDIR /app/notification-service

RUN go mod download

# Build the Go app
RUN go build -o notification-service .

# Expose the ports that the app will run on
EXPOSE 8082
EXPOSE 50052

# Command to run the executable
CMD ["./notification-service"]

