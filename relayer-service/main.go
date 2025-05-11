package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"relayer-service/database"
	"relayer-service/rabbitmq"
)

// RelayerEvent cấu trúc dữ liệu của event nhận từ notification-service
type RelayerEvent struct {
	EventType       string          `json:"event_type"`
	TransactionHash string          `json:"transaction_hash"`
	Depositor       string          `json:"depositor,omitempty"`
	Commitment      string          `json:"commitment,omitempty"`
	Recipient       string          `json:"recipient,omitempty"`
	NullifierHash   string          `json:"nullifier_hash,omitempty"`
	Data            json.RawMessage `json:"data"`
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
	queueName := "relayer_events"

	// Xử lý khi chương trình bị kết thúc
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Initialize the consumer
	consumer := &rabbitmq.Consumer{
		Queue: queueName,
	}
	err = consumer.Connect(rabbitmqURL)
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ consumer: %v", err)
	}
	defer consumer.Close()

	// Start consuming messages in a goroutine
	go func() {
		log.Println("Relayer Service is consuming messages from relayer_events...")
		err = consumer.Consume(func(message string) {
			log.Printf("Consumed message: %s", message)

			// Parse message để lấy thông tin event
			var event RelayerEvent
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
