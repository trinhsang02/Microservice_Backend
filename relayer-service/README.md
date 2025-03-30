# Relayer Service

The Relayer Service is a microservice responsible for managing the relaying of messages between different components of the system. It listens for incoming messages and forwards them to the appropriate destination.

## Setup Instructions

1. **Clone the Repository**
   ```bash
   git clone https://github.com/yourusername/golang-microservices.git
   cd golang-microservices/relayer-service
   ```

2. **Install Dependencies**
   Ensure you have Go installed on your machine. Run the following command to download the necessary dependencies:
   ```bash
   go mod tidy
   ```

3. **Configuration**
   Update any necessary configuration settings in the `main.go` file or through environment variables as required by your application.

4. **Run the Service**
   You can run the Relayer Service using the following command:
   ```bash
   go run main.go
   ```

## Usage

The Relayer Service will start listening for messages. Ensure that the other services (Blockchain Listener and Notification Service) are running and properly configured to communicate with the Relayer Service.

## API Endpoints

- **POST /relay**
  - Description: Relays a message to the specified destination.
  - Request Body: JSON object containing the message details.

## Contributing

If you would like to contribute to the Relayer Service, please fork the repository and submit a pull request with your changes.

## License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.