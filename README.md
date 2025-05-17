  # Golang Microservices Project

This project is a collection of microservices built using Go (Golang). It includes three main services:

1. **Blockchain Listener**: Listens for blockchain events and publishes them to RabbitMQ.
2. **Notification Service**: Consumes blockchain events from RabbitMQ, processes them, and publishes notifications to another queue.
3. **Relayer Service**: Consumes notifications from RabbitMQ and processes them further.

## Project Structure

```
golang-microservices
├── blockchain-listener
│   ├── main.go
│   ├── go.mod
│   ├── go.sum
│   ├── Dockerfile
│   ├── rabbitmq/
│   │   ├── producer.go
│   │   └── consumer.go
│   └── README.md
├── notification-service
│   ├── main.go
│   ├── go.mod
│   ├── go.sum
│   ├── Dockerfile

│   ├── rabbitmq/
│   │   ├── producer.go
│   │   └── consumer.go
│   └── README.md
├── relayer-service
│   ├── main.go
│   ├── go.mod
│   ├── go.sum
│   ├── Dockerfile
│   ├── rabbitmq/
│   │   ├── producer.go
│   │   └── consumer.go
│   └── README.md
├── docker-compose.yml
└── README.md
```

## Getting Started

### Prerequisites

- Go 1.22
- Docker and Docker Compose

### Installation

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd golang-microservices
   ```

2. Navigate to each service directory and install dependencies:
   ```bash
   cd blockchain-listener
   go mod tidy
   cd ../notification-service
   go mod tidy
   cd ../relayer-service
   go mod tidy
   ```

### Running the Services

You can use Docker Compose to run all services together. From the root of the project, execute:

```bash
docker-compose up
```

This will start the following services:
- **RabbitMQ**: Message broker for communication between services.
- **Blockchain Listener**: Publishes blockchain events to RabbitMQ.
- **Notification Service**: Consumes blockchain events and publishes notifications.
- **Relayer Service**: Consumes notifications and processes them.

### Environment Variables

Each service uses the following environment variables:
- `RABBITMQ_URL`: The RabbitMQ connection URL (default: `amqp://guest:guest@rabbitmq:5672/`).

### Service Details

#### **Blockchain Listener**
- **Description**: Publishes blockchain events to the `blockchain_events` queue.
- **Ports**:
  - `8081`: HTTP port (if applicable).
  - `50051`: gRPC port (if applicable).

#### **Notification Service**
- **Description**: Consumes events from the `blockchain_events` queue, processes them, and publishes notifications to the `relayer_events` queue.
- **Ports**:
  - `8082`: HTTP port (if applicable).
  - `50052`: gRPC port (if applicable).

#### **Relayer Service**
- **Description**: Consumes notifications from the `relayer_events` queue and processes them further.
- **Ports**:
  - `8083`: HTTP port (if applicable).
  - `50053`: gRPC port (if applicable).

## Testing

To test the services:
1. Start the services using Docker Compose.
2. Use tools like `curl`, `Postman`, or `grpcurl` to interact with the services.
3. Monitor the logs to verify message flow between services.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.

## License

This project is licensed under the MIT License. See the LICENSE file for details.