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

// SendMessage wraps the given message string in an interfaces.Message and sends
// it to the underlying message queue.
//
// A random UUID is generated and assigned to the Message's CorrelationId. This
// UUID is returned.
func (c *channel) SendMessage(message string) (correlationId string, err error) {
	correlationId = uuid.New().String()
	err = c.SendResponse(message, correlationId)

	return correlationId, err
}

// SendResponse wraps the given message string in an interfaces.Message and sends
// it to the underlying message queue.
//
// The given correlationId is assigned the the Message's CorrelationId.
func (c *channel) SendResponse(message string, correlationId string) error {
	return c.queueChannel.Publish("", c.queueName, false, false, interfaces.Message{
		ContentType:   "text/plain",
		CorrelationId: correlationId,
		Body:          []byte(message),
	})
}
