# Use an official Golang runtime as a parent image
FROM golang:1.24.1

# Set the working directory for the build
WORKDIR /app

# Copy only the needed folders
COPY blockchain-listener/ ./blockchain-listener/
COPY db/ ./db/
COPY mq/ ./mq/

# Set the working directory to the service
WORKDIR /app/blockchain-listener

# Build the service
RUN go build -o blockchain-listener .

# Expose the ports that the app will run on
EXPOSE 8081
EXPOSE 50051

# Command to run the executable
CMD ["./blockchain-listener"]