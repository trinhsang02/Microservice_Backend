package main

import (
	"log"
	"os"

	"relayer-service/rabbitmq"
)

func main() {
	// RabbitMQ connection details
	rabbitmqURL := os.Getenv("RABBITMQ_URL")
	if rabbitmqURL == "" {
		rabbitmqURL = "amqp://guest:guest@rabbitmq:5672/"
	}
	queueName := "relayer_events"

	// Initialize the consumer
	consumer := &rabbitmq.Consumer{
		Queue: queueName,
	}
	err := consumer.Connect(rabbitmqURL)
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ consumer: %v", err)
	}
	defer consumer.Close()
	err = consumer.Consume(func(message string) {
		log.Printf("Consumed message: %s", message)
		// Process the message here
	})
	if err != nil {
		log.Fatalf("Failed to consume messages: %v", err) // Use log.Fatalf to terminate with an error message
	}
	// Start consuming messages
	log.Println("Relayer Service is consuming messages from relayer_events...")
	err = consumer.Consume(func(message string) {
		log.Printf("Consumed message: %s", message)
		// Process the message here
	})
	if err != nil {
		log.Printf("Failed to consume messages: %v", err)
	}

	// Prevent the service from exiting
	select {}
}
