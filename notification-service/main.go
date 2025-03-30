package main

import (
	"log"
	"os"

	"notification-service/rabbitmq"
)

func main() {
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
	err := consumer.Connect()
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

	// Start consuming messages
	log.Println("Notification Service is consuming messages from blockchain_events...")
	err = consumer.Consume(func(message string) {
		log.Printf("Consumed message: %s", message)
		// Process the message and publish to relayer_events
		err := producer.Publish("Relayer event processed: " + message)
		if err != nil {
			log.Printf("Failed to publish message to relayer_events: %v", err)
		}
	})
	if err != nil {
		log.Printf("Failed to consume messages: %v", err)
	}
}
