# Use an official Golang runtime as a parent image
FROM golang:1.23

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o blockchain-listener .

# Expose the ports that the app will run on
EXPOSE 8081
EXPOSE 50051

# Command to run the executable
CMD ["./blockchain-listener"]