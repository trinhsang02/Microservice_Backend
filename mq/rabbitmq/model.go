package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type MQMessage struct {
    Type string      `json:"type"` // routingKey, ex: "deposit" or "withdrawal"
    Data interface{} `json:"data"`
}

type Producer struct {
    conn        *amqp.Connection
    channel     *amqp.Channel
    exchange    string
    exchangeType string
}

type Consumer struct {
    conn        *amqp.Connection
    channel     *amqp.Channel
    queueName   string
    exchange    string
    exchangeType string
    routingKeys []string
}
