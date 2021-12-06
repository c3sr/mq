package mq

import (
	"github.com/c3sr/mq/interfaces"
	"github.com/google/uuid"
)

// channel implements interfaces.Channel.
type channel struct {
	queueName    string
	queueChannel interfaces.QueueChannel
}

// SendMessage sends the given interfaces.Message to the underlying message queue.
//
// A random UUID is generated and assigned to the message's CorrelationId. This
// UUID is returned.
func (c *channel) SendMessage(message string) (string, error) {
	correlationId := uuid.New().String()
	err := c.queueChannel.Publish("", c.queueName, false, false, interfaces.Message{
		ContentType:   "text/plain",
		CorrelationId: correlationId,
		Body:          []byte(message),
	})

	return correlationId, err
}
