module relayer-service

go 1.24.1

require github.com/yourusername/yourrepo/mq v0.0.0

require github.com/rabbitmq/amqp091-go v1.10.0 // indirect

replace github.com/yourusername/yourrepo/mq => ../mq
