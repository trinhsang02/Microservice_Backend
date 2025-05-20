package main

import (
	"log"
	"os"

	"github.com/yourusername/yourrepo/mq/rabbitmq"
)

func main() {
	// RabbitMQ connection details
	rabbitmqURL := os.Getenv("RABBITMQ_URL")
	if rabbitmqURL == "" {
		rabbitmqURL = "amqp://guest:guest@rabbitmq:5672/"
	}
	queueName := "relayer_events"

	log.Printf("Connecting to RabbitMQ at %s", rabbitmqURL)
	log.Printf("Using queue: %s", queueName)

	// Initialize the consumer
	consumer, err := rabbitmq.NewConsumer(rabbitmqURL, "relayer_exchange", "topic", queueName, []string{"relayer.*"})
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ consumer: %v", err)
	}
	defer consumer.Close()

	// Start consuming messages
	log.Println("Relayer Service is consuming messages from relayer_events...")
	err = consumer.Consume(func(message rabbitmq.MQMessage) {
		log.Printf("Received message type: %s", message.Type)
		log.Printf("Received message data: %+v", message.Data)
		// Process the message here
	})
	if err != nil {
		log.Printf("Failed to consume messages: %v", err)
	}

	// Prevent the service from exiting
	select {}
}
