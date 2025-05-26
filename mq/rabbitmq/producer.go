package rabbitmq

import (
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

func NewProducer(url, exchange, exchangeType string) (*Producer, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

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
	return &Producer{
		conn:         conn,
		channel:      ch,
		exchange:     exchange,
		exchangeType: exchangeType,
	}, nil
}

func (p *Producer) PublishStruct(routingKey string, data interface{}) error {
	msg := MQMessage{
		Type: routingKey,
		Data: data,
	}
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return p.channel.Publish(
		p.exchange,
		routingKey,
		true,  // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (p *Producer) Close() {
	if p.channel != nil {
		_ = p.channel.Close()
	}
	if p.conn != nil {
		_ = p.conn.Close()
	}
}
