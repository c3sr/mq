package mq

import "github.com/c3sr/mq/interfaces"

type channel struct {
	exchangeName string
	queueChannel interfaces.QueueChannel
}

func (c *channel) SendMessage(message string) error {
	return c.queueChannel.Publish(c.exchangeName, "", false, false, interfaces.Message{
		ContentType: "text/plain",
		Body:        []byte(message),
	})
}
