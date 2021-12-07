// Package interfaces defines the interfaces used to interact with a message queue
package interfaces

// MessageQueue defines interactions available for a message queue implementation
type MessageQueue interface {
	Acknowledge(message Message) error
	Nack(message Message) error
	Shutdown()
	GetPublishChannel(name string) (Channel, error)
	SubscribeToChannel(name string) (<-chan Message, error)
}

// Message defines the structure of a message sent to or received from a message queue
type Message struct {
	Body          []byte
	ContentType   string
	CorrelationId string
	Id            uint64
}

// Queue defines the known details of a queue.
type Queue struct {
	Name string
}

// QueueChannel is the interface used to interact with an underlying message queue channel implementation.
type QueueChannel interface {
	Ack(tag uint64, multiple bool) error
	Close() error
	Consume(queue string, consumer string, autoAck bool, exclusive bool, noLocal bool, noWait bool, args map[string]interface{}) (<-chan Message, error)
	ExchangeDeclare(name string, kind string, durable bool, autoDelete bool, internal bool, noWait bool, args map[string]interface{}) error
	Nack(tag uint64, multiple bool, requeue bool) error
	Publish(exchange string, routingKey string, mandatory bool, immediate bool, message Message) error
	QueueBind(name string, routingKey string, exchange string, noWait bool, args map[string]interface{}) error
	QueueDeclare(name string, durable bool, autoDelete bool, exclusive bool, noWait bool, args map[string]interface{}) (Queue, error)
}

// QueueConnection is the interface used to interact with an underlying message queue connection implementation.
type QueueConnection interface {
	Channel() (QueueChannel, error)
	Close() error
}

// QueueDialer is the interface used to connect to an underlying message queue implementation.
type QueueDialer interface {
	Dial() (QueueConnection, error)
	URL() string
}

// Channel defines the interface for sending messages to an underlying message queue implementation.
type Channel interface {
	SendMessage(message string) (string, error)
}
