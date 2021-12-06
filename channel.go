package mq

import "github.com/c3sr/mq/interfaces"

type channel struct {
	queueName    string
	queueChannel interfaces.QueueChannel
}

func (c *channel) SendMessage(message string) error {
	return c.queueChannel.Publish("", c.queueName, false, false, interfaces.Message{
		ContentType: "text/plain",
		Body:        []byte(message),
	})
}
