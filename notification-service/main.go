package main

import (
	"log"
	"time"

	"github.com/yourusername/yourrepo/mq/rabbitmq"
)

type DepositNotification struct {
	DepositID string `json:"deposit_id"`
	Amount    string `json:"amount"`
	From      string `json:"from"`
	To        string `json:"to"`
}

func main() {
	// RabbitMQ connection details
	log.Println("Connecting to RabbitMQ...")
	producer, err := rabbitmq.NewProducer("amqp://guest:guest@rabbitmq:5672/", "relayer_exchange", "topic")
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ producer: %v", err)
	}
	defer producer.Close()
	log.Println("Successfully connected to RabbitMQ")

	// Create a notification
	notification := DepositNotification{
		DepositID: "123",
		Amount:    "100",
		From:      "0x123",
		To:        "0x456",
	}

	// Publish the notification
	log.Printf("Publishing notification to exchange 'relayer_exchange' with routing key 'relayer.deposit'")
	err = producer.PublishStruct("relayer.deposit", notification)
	if err != nil {
		log.Fatalf("Failed to publish notification: %v", err)
	}

	log.Println("Successfully published notification")

	// Keep the service running
	for {
		time.Sleep(time.Second)
	}
}
