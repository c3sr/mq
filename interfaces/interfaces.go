package interfaces

type MessageQueue interface {
	Acknowledge(message Message) error
	Shutdown()
	GetPublishChannel(name string) (Channel, error)
	SubscribeToChannel(name string) (<-chan Message, error)
}

type Message struct {
	Body          []byte
	ContentType   string
	CorrelationId string
	Id            uint64
}

type Queue struct {
	Name string
}

type QueueChannel interface {
	Ack(tag uint64, multiple bool) error
	Close() error
	Consume(queue string, consumer string, autoAck bool, exclusive bool, noLocal bool, noWait bool, args map[string]interface{}) (<-chan Message, error)
	ExchangeDeclare(name string, kind string, durable bool, autoDelete bool, internal bool, noWait bool, args map[string]interface{}) error
	Publish(exchange string, routingKey string, mandatory bool, immediate bool, message Message) error
	QueueBind(name string, routingKey string, exchange string, noWait bool, args map[string]interface{}) error
	QueueDeclare(name string, durable bool, autoDelete bool, exclusive bool, noWait bool, args map[string]interface{}) (Queue, error)
}

type QueueConnection interface {
	Channel() (QueueChannel, error)
	Close() error
}

type QueueDialer interface {
	Dial(url string) (QueueConnection, error)
}

type Channel interface {
	SendMessage(message string) (string, error)
}
