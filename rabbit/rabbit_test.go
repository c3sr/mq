// +build integration

package rabbit

import (
	"github.com/c3sr/mq"
	"github.com/stretchr/testify/assert"
	"testing"
)

func setupDialer() {
	dialer := NewRabbitDialer()
	mq.SetDialer(dialer)
}

func TestRabbitBasicCommunication(t *testing.T) {
	setupDialer()
	messageQueue, err := mq.NewMessageQueue()

	assert.Nil(t, err, "NewMessageQueue should not return an error")

	channel, err := messageQueue.GetPublishChannel("integration")

	assert.Nil(t, err, "GetPublishChannel should not return an error")

	messages, err := messageQueue.SubscribeToChannel("integration")

	assert.Nil(t, err, "SubscribeToChannel should not return an error")

	channel.SendMessage("test message")
	message := <-messages

	assert.Equal(t, "test message", string(message.Body))

	messageQueue.Shutdown()
}

func TestMessageAcknowledgement(t *testing.T) {
	setupDialer()
	messageQueue, _ := mq.NewMessageQueue()
	messageQueue2, _ := mq.NewMessageQueue()
	messageQueue3, _ := mq.NewMessageQueue()

	channel, _ := messageQueue.GetPublishChannel("acknowledgement")
	messages2, _ := messageQueue2.SubscribeToChannel("acknowledgement")
	messages3, _ := messageQueue3.SubscribeToChannel("acknowledgement")

	channel.SendMessage("acknowledged")

	message2 := <-messages2
	err := messageQueue2.Acknowledge(message2)

	if err != nil {
		t.Errorf("failed to acknowledge message: %s", err)
	}

	select {
	case msg := <-messages3:
		t.Errorf("Should not have received acknowledged message: %s", string(msg.Body))

	default:
		println("Acknowledge successful")
	}

	messageQueue.Shutdown()
	messageQueue2.Shutdown()
	messageQueue3.Shutdown()
}

func TestCorrelationIdPropagation(t *testing.T) {
	setupDialer()
	messageQueue, _ := mq.NewMessageQueue()
	messageQueue2, _ := mq.NewMessageQueue()

	channel, _ := messageQueue.GetPublishChannel("correlation")
	messages, _ := messageQueue2.SubscribeToChannel("correlation")

	correlationId, _ := channel.SendMessage("correlated message")
	message := <-messages

	assert.Equal(t, correlationId, message.CorrelationId)

	messageQueue.Shutdown()
	messageQueue2.Shutdown()
}

func TestMessageNegativeAcknowledgement(t *testing.T) {
	setupDialer()
	messageQueue, _ := mq.NewMessageQueue()
	messageQueue2, _ := mq.NewMessageQueue()
	messageQueue3, _ := mq.NewMessageQueue()

	channel, _ := messageQueue.GetPublishChannel("nacknowledgment")
	messages2, _ := messageQueue2.SubscribeToChannel("nacknowledgment")
	messages3, _ := messageQueue3.SubscribeToChannel("nacknowledgment")

	correlationId, _ := channel.SendMessage("nacked")

	message2 := <-messages2
	err := messageQueue2.Nack(message2)

	if err != nil {
		t.Errorf("failed to nack message: %s", err)
	}

	message3 := <-messages3

	assert.Equal(t, correlationId, message3.CorrelationId)

	messageQueue.Shutdown()
	messageQueue2.Shutdown()
	messageQueue3.Shutdown()
}
