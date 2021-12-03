// +build integration

package rabbit

import (
	"github.com/c3sr/mq"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRabbitBasicCommunication(t *testing.T) {
	dialer := NewRabbitDialer()
	mq.SetDialer(dialer)
	messageQueue, err := mq.NewMessageQueue()

	assert.Nil(t, err, "NewMessageQueue should not return an error")

	channel, err := messageQueue.GetPublishChannel("integration")

	assert.Nil(t, err, "GetPublishChannel should not return an error")

	messages, err := messageQueue.SubscribeToChannel("integration")

	assert.Nil(t, err, "SubscribeToChannel should not return an error")

	channel.SendMessage("test message")
	message := <-messages

	assert.Equal(t, "test message", string(message.Body))
}
