package interfaces

type MessageQueue interface {
	Shutdown()
	GetPublishChannel(name string) (Channel, error)
	SubscribeToChannel(name string) (<-chan Message, error)
}

type Message struct {
	ContentType string
	Body        []byte
}

type Queue struct {
	Name string
}

type QueueChannel interface {
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
	SendMessage(message string) error
}
