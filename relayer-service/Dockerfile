# Use an official Golang runtime as a parent image
FROM golang:1.24

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files first
COPY relayer-service/ ./relayer-service/

# Copy the mq module
COPY mq/ ./mq/

WORKDIR /app/relayer-service
# Download dependencies
RUN go mod download

# Build the Go app
RUN go build -o relayer-service .

# Expose the ports that the app will run on
EXPOSE 8083
EXPOSE 50053

# Command to run the executable
CMD ["./relayer-service"]