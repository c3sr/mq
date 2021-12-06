package mq

import (
	"github.com/c3sr/mq/interfaces"
)

type messageQueue struct {
	channel    interfaces.QueueChannel
	connection interfaces.QueueConnection
}

func (q *messageQueue) Acknowledge(message interfaces.Message) error {
	return q.channel.Ack(message.Id, false)
}

func (q *messageQueue) Nack(message interfaces.Message) error {
	return q.channel.Nack(message.Id, false, true)
}

func (q *messageQueue) Shutdown() {
	q.channel.Close()
	q.connection.Close()
}

func (q *messageQueue) GetPublishChannel(name string) (interfaces.Channel, error) {
	q.channel.QueueDeclare(
		name,
		false,
		false,
		false,
		false,
		nil)

	return &channel{
		queueName:    name,
		queueChannel: q.channel,
	}, nil
}

func (q *messageQueue) SubscribeToChannel(name string) (<-chan interfaces.Message, error) {
	queue, err := q.channel.QueueDeclare(name, false, false, false, false, nil)

	if err != nil {
		return nil, err
	}

	messages, err := q.channel.Consume(queue.Name, "", false, false, false, false, nil)

	return messages, err
}
