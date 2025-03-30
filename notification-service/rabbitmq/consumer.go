package rabbitmq

import (
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	Queue    string
	conn     *amqp.Connection
	channel  *amqp.Channel
	messages <-chan amqp.Delivery
}

// Connect establishes a connection to RabbitMQ and initializes the consumer.
func (c *Consumer) Connect() error {
    var err error

    // Get the RabbitMQ URL from the environment variable or use a default value
    rabbitmqURL := os.Getenv("RABBITMQ_URL")
    if rabbitmqURL == "" {
        rabbitmqURL = "amqp://guest:guest@localhost:5672/" // Default RabbitMQ URL
    }

    // Establish a connection to RabbitMQ
    c.conn, err = amqp.Dial(rabbitmqURL)
    if err != nil {
        return err
    }

    // Open a channel
    c.channel, err = c.conn.Channel()
    if err != nil {
        return err
    }

    // Start consuming messages from the queue
    c.messages, err = c.channel.Consume(
        c.Queue, // queue
        "",      // consumer
        true,    // auto-ack
        false,   // exclusive
        false,   // no-local
        false,   // no-wait
        nil,     // args
    )
    return err
}

// Close closes the RabbitMQ connection and channel.
func (c *Consumer) Close() {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
}

// Consume processes messages from the queue using the provided handler function.
func (c *Consumer) Consume(handler func(string)) error {
	for msg := range c.messages {
		handler(string(msg.Body))
	}
	return nil
}
