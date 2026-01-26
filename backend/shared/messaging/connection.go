package messaging

// Connection represents a RabbitMQ connection
type Connection interface {
	Close() error
}
