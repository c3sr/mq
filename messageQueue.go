package mq

import (
	"github.com/c3sr/mq/interfaces"
	amqp "github.com/rabbitmq/amqp091-go"
)

type messageQueue struct {
	channel interfaces.QueueChannel
}

func (q messageQueue) Acknowledge(message interfaces.Message) error {
	return qChannel.Ack(message.Id, false)
}

func (q messageQueue) Shutdown() {
	qChannel.Close()
	qChannel = nil
	connection.Close()
	connection = nil
}

func (q messageQueue) GetPublishChannel(name string) (interfaces.Channel, error) {
	q.channel.ExchangeDeclare(
		name,
		amqp.ExchangeDirect,
		false,
		false,
		false,
		false,
		nil)

	return &channel{
		exchangeName: name,
		queueChannel: q.channel,
	}, nil
}

func (q messageQueue) SubscribeToChannel(name string) (<-chan interfaces.Message, error) {
	err := q.channel.ExchangeDeclare(
		name,
		amqp.ExchangeDirect,
		false,
		false,
		false,
		false,
		nil)

	if err != nil {
		return nil, err
	}

	queue, err := q.channel.QueueDeclare("", false, false, true, false, nil)

	if err != nil {
		return nil, err
	}

	err = q.channel.QueueBind(queue.Name, "", name, false, nil)

	if err != nil {
		return nil, err
	}

	messages, err := q.channel.Consume(queue.Name, "", false, false, false, false, nil)

	return messages, err
}
