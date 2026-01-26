package messaging

// Publisher handles publishing events to RabbitMQ
type Publisher interface {
	Publish(eventType string, event interface{}) error
}
