# Use an official Golang runtime as a parent image
FROM golang:1.24

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the source code into the container
COPY common-service/ ./common-service/
COPY db/ ./db/

# Set working directory to the service
WORKDIR /app/common-service

# Build the Go app
RUN go build -o common-service .

# Expose the ports that the app will run on
EXPOSE 8080

# Command to run the executable
CMD ["./common-service"]

