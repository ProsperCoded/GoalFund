package messaging

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gofund/shared/metrics"
	"github.com/streadway/amqp"
)

// RabbitMQConnection represents a RabbitMQ connection
type RabbitMQConnection struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

// NewRabbitMQConnection creates a new RabbitMQ connection
func NewRabbitMQConnection(url string) (*RabbitMQConnection, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	return &RabbitMQConnection{
		conn:    conn,
		channel: ch,
	}, nil
}

// Close closes the RabbitMQ connection
func (r *RabbitMQConnection) Close() error {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		return r.conn.Close()
	}
	return nil
}

// RabbitMQPublisher implements the Publisher interface
type RabbitMQPublisher struct {
	channel      *amqp.Channel
	exchangeName string
}

// NewRabbitMQPublisher creates a new RabbitMQ publisher
func NewRabbitMQPublisher(conn *RabbitMQConnection, exchangeName string) (*RabbitMQPublisher, error) {
	// Declare exchange
	err := conn.channel.ExchangeDeclare(
		exchangeName, // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	return &RabbitMQPublisher{
		channel:      conn.channel,
		exchangeName: exchangeName,
	}, nil
}

// Publish publishes an event to RabbitMQ with Datadog metrics
func (p *RabbitMQPublisher) Publish(eventType string, event interface{}) error {
	start := time.Now()
	
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	routingKey := fmt.Sprintf("events.%s", eventType)

	err = p.channel.Publish(
		p.exchangeName, // exchange
		routingKey,     // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
			Timestamp:   time.Now(),
		},
	)
	
	duration := time.Since(start)
	
	// Track event publishing metrics
	if err != nil {
		metrics.TrackEventPublished(eventType, false, duration)
		return fmt.Errorf("failed to publish event: %w", err)
	}
	
	metrics.TrackEventPublished(eventType, true, duration)
	log.Printf("Published event: %s (duration: %v)", eventType, duration)
	return nil
}

// RabbitMQConsumer implements the Consumer interface
type RabbitMQConsumer struct {
	channel      *amqp.Channel
	exchangeName string
	queueName    string
}

// NewRabbitMQConsumer creates a new RabbitMQ consumer
func NewRabbitMQConsumer(conn *RabbitMQConnection, exchangeName, queueName string) (*RabbitMQConsumer, error) {
	// Declare exchange
	err := conn.channel.ExchangeDeclare(
		exchangeName, // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	// Declare queue
	_, err = conn.channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	return &RabbitMQConsumer{
		channel:      conn.channel,
		exchangeName: exchangeName,
		queueName:    queueName,
	}, nil
}

// Consume consumes events from RabbitMQ with Datadog metrics
func (c *RabbitMQConsumer) Consume(eventType string, handler func([]byte) error) error {
	routingKey := fmt.Sprintf("events.%s", eventType)

	// Bind queue to exchange
	err := c.channel.QueueBind(
		c.queueName,    // queue name
		routingKey,     // routing key
		c.exchangeName, // exchange
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to bind queue: %w", err)
	}

	msgs, err := c.channel.Consume(
		c.queueName, // queue
		"",          // consumer
		true,        // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	go func() {
		for msg := range msgs {
			start := time.Now()
			
			// Calculate event age (time since published)
			var eventAge time.Duration
			if !msg.Timestamp.IsZero() {
				eventAge = time.Since(msg.Timestamp)
			}
			
			// Handle the message
			err := handler(msg.Body)
			processingDuration := time.Since(start)
			
			// Track event consumption metrics
			if err != nil {
				log.Printf("Error handling message: %v", err)
				metrics.TrackEventConsumed(eventType, false, processingDuration, eventAge)
			} else {
				metrics.TrackEventConsumed(eventType, true, processingDuration, eventAge)
			}
		}
	}()

	log.Printf("Started consuming events: %s", eventType)
	return nil
}