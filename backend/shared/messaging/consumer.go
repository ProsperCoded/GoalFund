package messaging

// Consumer handles consuming events from RabbitMQ
type Consumer interface {
	Consume(eventType string, handler func([]byte) error) error
}
