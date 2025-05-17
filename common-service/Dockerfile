# Build stage
FROM golang:1.23 

WORKDIR /app

# Copy only the needed folders
COPY common-service/ ./common-service/
COPY db/ ./db/

# Set the working directory to the service
WORKDIR /app/common-service

# Build the service
RUN go build -o common-service .
# Expose port
EXPOSE 8080

# Run the application
CMD ["./common-service"]

