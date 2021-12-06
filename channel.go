package mq

import (
	"github.com/c3sr/mq/interfaces"
	"github.com/google/uuid"
)

type channel struct {
	queueName    string
	queueChannel interfaces.QueueChannel
}

func (c *channel) SendMessage(message string) (string, error) {
	correlationId := uuid.New().String()
	err := c.queueChannel.Publish("", c.queueName, false, false, interfaces.Message{
		ContentType:   "text/plain",
		CorrelationId: correlationId,
		Body:          []byte(message),
	})

	return correlationId, err
}
