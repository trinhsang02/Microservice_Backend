package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
    Queue string
    conn  *amqp.Connection
    ch    *amqp.Channel
}

// Connect establishes a connection to RabbitMQ and declares the queue.
func (c *Consumer) Connect(rabbitmqURL string) error {
    var err error
    c.conn, err = amqp.Dial(rabbitmqURL)
    if err != nil {
        return err
    }

    c.ch, err = c.conn.Channel()
    if err != nil {
        return err
    }

    _, err = c.ch.QueueDeclare(
        c.Queue, // name
        true,    // durable
        false,   // delete when unused
        false,   // exclusive
        false,   // no-wait
        nil,     // arguments
    )
    return err
}

// Close closes the RabbitMQ connection and channel.
func (c *Consumer) Close() {
    if c.ch != nil {
        c.ch.Close()
    }
    if c.conn != nil {
        c.conn.Close()
    }
}

// Consume starts consuming messages from the queue.
func (c *Consumer) Consume(handler func(string)) error {
    msgs, err := c.ch.Consume(
        c.Queue, // queue
        "",      // consumer
        true,    // auto-ack
        false,   // exclusive
        false,   // no-local
        false,   // no-wait
        nil,     // args
    )
    if err != nil {
        return err
    }

    go func() {
        for d := range msgs {
            handler(string(d.Body))
        }
    }()

    return nil
}