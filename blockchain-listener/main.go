package main

import (
	"log"
	"os"
	"time"

	"blockchain-listener/rabbitmq"
)

func main() {
    // RabbitMQ connection details
    rabbitmqURL  := os.Getenv("RABBITMQ_URL")
    if rabbitmqURL == "" {
        rabbitmqURL = "amqp://guest:guest@rabbitmq:5672/"
    }
    queueName := "blockchain_events"

    // Initialize the producer
    producer, err := rabbitmq.NewProducer(rabbitmqURL , queueName)
    if err != nil {
        log.Fatalf("Failed to initialize RabbitMQ producer: %v", err)
    }
    defer producer.Close()

    // Initialize the consumer
    // consumer := &rabbitmq.Consumer{
    //     Queue: queueName,
    // }
    // err = consumer.Connect(rabbitmqURL)
    // if err != nil {
    //     log.Fatalf("Failed to initialize RabbitMQ consumer: %v", err)
    // }
    // defer consumer.Close()

    // Start consuming messages in a separate goroutine
    // go func() {
    //     log.Println("Blockchain Listener is consuming messages...")
    //     err := consumer.Consume(func(message string) {
    //         log.Printf("Consumed message: %s", message)
    //         // Process the message here
    //     })
    //     if err != nil {
    //         log.Printf("Failed to consume messages: %v", err)
    //     }
    // }()

    // Publish messages periodically
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()

    log.Println("Blockchain Listener is publishing messages...")
    for t := range ticker.C {
        message := "Blockchain event at " + t.String()
        err := producer.Publish(message)
        if err != nil {
            log.Printf("Failed to publish message: %v", err)
        }
    }
    
}