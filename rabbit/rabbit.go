package rabbit

import (
	"github.com/c3sr/mq/interfaces"
	amqp "github.com/rabbitmq/amqp091-go"
)

type rabbitChannel struct {
	*amqp.Channel
}

func (c *rabbitChannel) Close() error {
	return c.Channel.Close()
}

func (c *rabbitChannel) Consume(queue string, consumer string, autoAck bool, exclusive bool, noLocal bool, noWait bool, args map[string]interface{}) (<-chan interfaces.Message, error) {
	delivery, err := c.Channel.Consume(queue, consumer, autoAck, exclusive, noLocal, noWait, args)
	messages := make(chan interfaces.Message)

	if err != nil {
		return nil, err
	}

	go func() {
		for d := range delivery {
			messages <- interfaces.Message{
				Id:            d.DeliveryTag,
				ContentType:   d.ContentType,
				CorrelationId: d.CorrelationId,
				Body:          d.Body,
			}
		}
	}()

	return messages, nil
}

func (c *rabbitChannel) ExchangeDeclare(name string, kind string, durable bool, autoDelete bool, internal bool, noWait bool, args map[string]interface{}) error {
	return c.Channel.ExchangeDeclare(name, kind, durable, autoDelete, internal, noWait, args)
}

func (c *rabbitChannel) Publish(exchange string, routingKey string, mandatory bool, immediate bool, message interfaces.Message) error {
	return c.Channel.Publish(
		exchange,
		routingKey,
		mandatory,
		immediate,
		amqp.Publishing{
			ContentType:   message.ContentType,
			CorrelationId: message.CorrelationId,
			Body:          message.Body,
		})
}

func (c *rabbitChannel) QueueBind(name string, routingKey string, exchange string, noWait bool, args map[string]interface{}) error {
	return c.Channel.QueueBind(name, routingKey, exchange, noWait, args)
}

func (c *rabbitChannel) QueueDeclare(name string, durable bool, autoDelete bool, exclusive bool, noWait bool, args map[string]interface{}) (queue interfaces.Queue, err error) {
	rabbitQueue, err := c.Channel.QueueDeclare(name, durable, autoDelete, exclusive, noWait, args)

	if err != nil {
		return
	}

	return interfaces.Queue{Name: rabbitQueue.Name}, nil
}

type rabbitConnection struct {
	*amqp.Connection
}

func (c *rabbitConnection) Channel() (interfaces.QueueChannel, error) {
	channel, err := c.Connection.Channel()

	return &rabbitChannel{channel}, err
}

func (c *rabbitConnection) Close() error {
	return c.Connection.Close()
}
