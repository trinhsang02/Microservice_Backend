package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"notification-service/database"
	"notification-service/rabbitmq"
)

// BlockchainEvent cấu trúc dữ liệu của event nhận từ blockchain-listener
type BlockchainEvent struct {
	EventType       string          `json:"event_type"`
	TransactionHash string          `json:"transaction_hash"`
	BlockNumber     int64           `json:"block_number"`
	EventData       json.RawMessage `json:"event_data"`
}

func main() {
	// Khởi tạo kết nối PostgreSQL
	err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDB()

	// RabbitMQ connection details
	rabbitmqURL := os.Getenv("RABBITMQ_URL")
	if rabbitmqURL == "" {
		rabbitmqURL = "amqp://guest:guest@rabbitmq:5672/"
	}
	consumerQueue := "blockchain_events"
	producerQueue := "relayer_events"

	// Initialize the consumer
	consumer := &rabbitmq.Consumer{
		Queue: consumerQueue,
	}
	err = consumer.Connect()
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ consumer: %v", err)
	}
	defer consumer.Close()

	// Initialize the producer
	producer, err := rabbitmq.NewProducer(rabbitmqURL, producerQueue)
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ producer: %v", err)
	}
	defer producer.Close()

	// Xử lý khi chương trình bị kết thúc
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start consuming messages in a goroutine
	go func() {
		// Start consuming messages
		log.Println("Notification Service is consuming messages from blockchain_events...")
		err = consumer.Consume(func(message string) {
			log.Printf("Consumed message: %s", message)

			// Parse message để lấy thông tin event
			var event BlockchainEvent
			err := json.Unmarshal([]byte(message), &event)
			if err != nil {
				log.Printf("Error parsing message: %v", err)
				return
			}

		})
		if err != nil {
			log.Printf("Failed to consume messages: %v", err)
		}
	}()

	// Chờ tín hiệu kết thúc
	<-sigChan
	log.Println("Shutting down gracefully...")
}


