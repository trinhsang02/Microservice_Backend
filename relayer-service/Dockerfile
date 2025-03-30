# Use an official Golang runtime as a parent image
FROM golang:1.22

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o relayer-service .

# Expose the ports that the app will run on
EXPOSE 8083
EXPOSE 50053

# Command to run the executable
CMD ["./relayer-service"]