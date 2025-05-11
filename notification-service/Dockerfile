# Use an official Golang runtime as a parent image
FROM golang:1.23
# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o notification-service .

# Expose the ports that the app will run on
EXPOSE 8082
EXPOSE 50052

# Command to run the executable
CMD ["./notification-service"]