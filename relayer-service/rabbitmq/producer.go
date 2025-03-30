package rabbitmq

import (
    amqp "github.com/rabbitmq/amqp091-go"
)

type Producer struct {
    connection *amqp.Connection
    channel    *amqp.Channel
    queueName  string
}

// NewProducer initializes a new RabbitMQ producer
func NewProducer(url, queueName string) (*Producer, error) {
    conn, err := amqp.Dial(url)
    if err != nil {
        return nil, err
    }

    ch, err := conn.Channel()
    if err != nil {
        conn.Close()
        return nil, err
    }

    _, err = ch.QueueDeclare(
        queueName, // name
        true,      // durable
        false,     // delete when unused
        false,     // exclusive
        false,     // no-wait
        nil,       // arguments
    )
    if err != nil {
        ch.Close()
        conn.Close()
        return nil, err
    }

    return &Producer{
        connection: conn,
        channel:    ch,
        queueName:  queueName,
    }, nil
}

// Publish sends a message to the RabbitMQ queue
func (p *Producer) Publish(message string) error {
    return p.channel.Publish(
        "",          // exchange
        p.queueName, // routing key
        false,       // mandatory
        false,       // immediate
        amqp.Publishing{
            ContentType: "text/plain",
            Body:        []byte(message),
        },
    )
}

// Close closes the RabbitMQ connection and channel
func (p *Producer) Close() {
    if p.channel != nil {
        p.channel.Close()
    }
    if p.connection != nil {
        p.connection.Close()
    }
}