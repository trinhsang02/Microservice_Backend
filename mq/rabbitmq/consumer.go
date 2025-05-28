package rabbitmq

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func NewConsumer(url, exchange, exchangeType, queueName string, routingKeys []string) (*Consumer, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}
	// Khai báo exchange
	err = ch.ExchangeDeclare(
		exchange,
		exchangeType, // "direct"
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}
	// Khai báo queue
	q, err := ch.QueueDeclare(
		queueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}
	// Bind queue với từng routing key
	for _, key := range routingKeys {
		err = ch.QueueBind(
			q.Name,
			key,
			exchange,
			false,
			nil,
		)
		if err != nil {
			ch.Close()
			conn.Close()
			return nil, err
		}
	}

	return &Consumer{
		conn:         conn,
		channel:      ch,
		queueName:    q.Name,
		exchange:     exchange,
		exchangeType: exchangeType,
		routingKeys:  routingKeys,
	}, nil
}

func (c *Consumer) Consume(handler func(MQMessage)) error {
	msgs, err := c.channel.Consume(
		c.queueName, // queue
		"",          // consumer
		false,       // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
	if err != nil {
		return err
	}
	go func() {
		for d := range msgs {
			var msg MQMessage
			if err := json.Unmarshal(d.Body, &msg); err != nil {
				log.Println("Error unmarshalling message:", err)
				d.Nack(false, false)
				continue
			}
			handler(msg)
			d.Ack(false)
		}
	}()
	return nil
}

func (c *Consumer) Close() {
	if c.channel != nil {
		_ = c.channel.Close()
	}
	if c.conn != nil {
		_ = c.conn.Close()
	}
}
 