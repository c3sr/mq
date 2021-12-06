package mq

import (
	"github.com/c3sr/mq/interfaces"
)

type messageQueue struct {
	channel    interfaces.QueueChannel
	connection interfaces.QueueConnection
}

// Acknowledge sends an acknowledgement for the given interfaces.Message. In most cases this results in the message
// being discarded by the underlying queueing system so that it will not be delivered to any further clients.
func (q *messageQueue) Acknowledge(message interfaces.Message) error {
	return q.channel.Ack(message.Id, false)
}

// Nack sends a negative acknowledgement for the given interfaces.Message. In most cases this results in the message
// being re-queued by the underlying queueing system for delivery to another client.
func (q *messageQueue) Nack(message interfaces.Message) error {
	return q.channel.Nack(message.Id, false, true)
}

// Shutdown closes any connection to the underlying message queue and clears associated resources.
func (q *messageQueue) Shutdown() {
	q.channel.Close()
	q.connection.Close()
}

// GetPublishChannel opens a channel for publishing with the given name and returns the associated
// interfaces.Channel object.
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

// SubscribeToChannel connects to the underlying message queue with the given name and returns
// a channel of interfaces.Message objects.
func (q *messageQueue) SubscribeToChannel(name string) (<-chan interfaces.Message, error) {
	queue, err := q.channel.QueueDeclare(name, false, false, false, false, nil)

	if err != nil {
		return nil, err
	}

	messages, err := q.channel.Consume(queue.Name, "", false, false, false, false, nil)

	return messages, err
}
