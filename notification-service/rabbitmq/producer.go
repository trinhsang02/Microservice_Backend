package rabbitmq

import (
    "log"
    "time"
    amqp "github.com/rabbitmq/amqp091-go"
)

type Producer struct {
    connection   *amqp.Connection
    channel      *amqp.Channel
    queueName    string
    url          string
    closeChan    chan *amqp.Error
    stopChan     chan bool
}

// NewProducer initializes a new RabbitMQ producer
func NewProducer(url, queueName string) (*Producer, error) {
    producer := &Producer{
        queueName: queueName,
        url:       url,
        stopChan:  make(chan bool),
    }
    
    err := producer.connect()
    if err != nil {
        return nil, err
    }
    
    // Start monitoring the connection
    go producer.monitorConnection()
    
    return producer, nil
}

// connect establishes the initial connection
func (p *Producer) connect() error {
    var err error
    p.connection, err = amqp.Dial(p.url)
    if err != nil {
        return err
    }
    
    p.channel, err = p.connection.Channel()
    if err != nil {
        p.connection.Close()
        return err
    }
    
    _, err = p.channel.QueueDeclare(
        p.queueName, // name
        true,        // durable
        false,       // delete when unused
        false,       // exclusive
        false,       // no-wait
        nil,         // arguments
    )
    if err != nil {
        p.channel.Close()
        p.connection.Close()
        return err
    }
    
    // Setup notification for connection closure
    p.closeChan = make(chan *amqp.Error)
    p.connection.NotifyClose(p.closeChan)
    
    log.Println("Producer successfully connected to RabbitMQ")
    return nil
}

// Publish sends a message to the RabbitMQ queue
func (p *Producer) Publish(message string) error {
    if p.channel == nil || p.connection == nil || p.connection.IsClosed() {
        // If we reach here, it means our automatic reconnection hasn't completed yet
        // Try a manual reconnect
        err := p.reconnect()
        if err != nil {
            return err
        }
    }
    
    return p.channel.Publish(
        "",           // exchange
        p.queueName,  // routing key
        false,        // mandatory
        false,        // immediate
        amqp.Publishing{
            ContentType:  "text/plain",
            Body:         []byte(message),
            DeliveryMode: amqp.Persistent, // Make messages persistent
        },
    )
}

// monitorConnection watches for closed connections and reconnects
func (p *Producer) monitorConnection() {
    for {
        select {
        case <-p.stopChan:
            return
        case err := <-p.closeChan:
            if err != nil {
                log.Printf("RabbitMQ connection closed: %v", err)
            }
            p.reconnect()
        }
    }
}

// reconnect tries to re-establish the connection and channel
func (p *Producer) reconnect() error {
    // Close old connections if they exist
    if p.channel != nil {
        p.channel.Close()
    }
    if p.connection != nil {
        p.connection.Close()
    }
    
    // Reconnect with backoff
    backoff := 1 * time.Second
    maxBackoff := 30 * time.Second
    
    for {
        select {
        case <-p.stopChan:
            return nil
        default:
            log.Printf("Producer attempting to reconnect to RabbitMQ in %v...", backoff)
            time.Sleep(backoff)
            
            err := p.connect()
            if err == nil {
                log.Println("Producer successfully reconnected to RabbitMQ")
                return nil
            }
            
            log.Printf("Producer failed to reconnect: %v", err)
            
            // Increase backoff for next attempt, up to maximum
            backoff *= 2
            if backoff > maxBackoff {
                backoff = maxBackoff
            }
        }
    }
}

// Close stops the producer and closes connections
func (p *Producer) Close() {
    p.stopChan <- true
    
    if p.channel != nil {
        p.channel.Close()
    }
    if p.connection != nil {
        p.connection.Close()
    }
}