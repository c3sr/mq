// +build !integration

package mq

import (
	"github.com/c3sr/mq/interfaces"
	"github.com/stretchr/testify/assert"
	"testing"
)

type spyChannel struct {
	ackedMessageId       uint64
	boundExchangeName    string
	closeCalled          bool
	declaredExchangeName string
	declaredQueueName    string
	lastCorrelationId    string
	lastExchangeName     string
	lastMessage          string
	lastRoutingKey       string
	nackedMessageId      uint64
	nackRequeue          bool
	consumedQueueName    string
}

func (s *spyChannel) Close() error {
	s.closeCalled = true

	return nil
}

func (s *spyChannel) Ack(tag uint64, multiple bool) error {
	s.ackedMessageId = tag
	return nil
}

func (s *spyChannel) Consume(queue string, consumer string, autoAck bool, exclusive bool, noLocal bool, noWait bool, args map[string]interface{}) (<-chan interfaces.Message, error) {
	messages := make(chan interfaces.Message)
	s.consumedQueueName = queue

	return messages, nil
}

func (s *spyChannel) ExchangeDeclare(name string, kind string, durable bool, autoDelete bool, internal bool, noWait bool, args map[string]interface{}) error {
	s.declaredExchangeName = name

	return nil
}

func (s *spyChannel) Nack(tag uint64, multiple bool, requeue bool) error {
	s.nackedMessageId = tag
	s.nackRequeue = requeue
	return nil
}

func (s *spyChannel) Publish(exchange string, routingKey string, mandatory bool, immediate bool, message interfaces.Message) error {
	s.lastExchangeName = exchange
	s.lastRoutingKey = routingKey
	s.lastMessage = string(message.Body)
	s.lastCorrelationId = message.CorrelationId

	return nil
}

func (s *spyChannel) QueueBind(name string, routingKey string, exchange string, noWait bool, args map[string]interface{}) error {
	s.boundExchangeName = exchange

	return nil
}

func (s *spyChannel) QueueDeclare(name string, durable bool, autoDelete bool, exclusive bool, noWait bool, args map[string]interface{}) (interfaces.Queue, error) {
	s.declaredQueueName = name
	queue := interfaces.Queue{Name: "spyQueue"}

	return queue, nil
}

type spyConnection struct {
	channel *spyChannel
}

func (s *spyConnection) Channel() (interfaces.QueueChannel, error) {
	s.channel = &spyChannel{
		closeCalled:          false,
		declaredExchangeName: "",
	}

	return s.channel, nil
}

func (s *spyConnection) Close() error {
	return nil
}

func (s *spyConnection) NotifyClose(chan error) {
}

type spyDialer struct {
	connection *spyConnection
}

func (d *spyDialer) Dial() (interfaces.QueueConnection, error) {
	d.connection = &spyConnection{}

	return d.connection, nil
}

func (d *spyDialer) URL() string {
	return ""
}

var testDialer *spyDialer

func setupDialer() {
	testDialer = &spyDialer{}

	SetDialer(testDialer)
}

func TestNoDialerError(t *testing.T) {
	_, err := NewMessageQueue()

	assert.Equal(t, "A dialer must be provided using SetDialer()", err.Error())
}

func TestGetPublishChannelDeclaresQueue(t *testing.T) {
	setupDialer()
	mq, _ := NewMessageQueue()

	channel, _ := mq.GetPublishChannel("test")

	assert.NotNil(t, channel)
	assert.Equal(t, "test", testDialer.connection.channel.declaredQueueName)
}

func TestSendMessagePublishesToQueueChannel(t *testing.T) {
	setupDialer()
	mq, _ := NewMessageQueue()

	channel, err := mq.GetPublishChannel("publish")
	_, err = channel.SendMessage("hello, world!")

	assert.Nil(t, err)
	assert.Equal(t, "", testDialer.connection.channel.lastExchangeName)
	assert.Equal(t, "publish", testDialer.connection.channel.lastRoutingKey)
	assert.Equal(t, "hello, world!", testDialer.connection.channel.lastMessage)
}

func TestSubscribeChannel(t *testing.T) {
	setupDialer()
	mq, _ := NewMessageQueue()

	messages, _ := mq.SubscribeToChannel("consume")

	assert.Equal(t, "", testDialer.connection.channel.declaredExchangeName)
	assert.Equal(t, "consume", testDialer.connection.channel.declaredQueueName)
	assert.Equal(t, "", testDialer.connection.channel.boundExchangeName)
	assert.Equal(t, "spyQueue", testDialer.connection.channel.consumedQueueName)
	assert.NotNil(t, messages)
}

func TestAcknowledgeMessage(t *testing.T) {
	setupDialer()
	mq, _ := NewMessageQueue()
	message := interfaces.Message{Id: 1}

	mq.Acknowledge(message)

	assert.Equal(t, uint64(1), testDialer.connection.channel.ackedMessageId)
}

func TestSendMessageSetsCorrelationIdToUUID(t *testing.T) {
	setupDialer()
	mq, _ := NewMessageQueue()
	channel, _ := mq.GetPublishChannel("publish")
	correlationId, _ := channel.SendMessage("test")

	assert.Equal(t, correlationId, testDialer.connection.channel.lastCorrelationId)
}

func TestSendResponseUsesPassedCorrelationId(t *testing.T) {
	setupDialer()
	mq, _ := NewMessageQueue()
	channel, _ := mq.GetPublishChannel("publish")
	_ = channel.SendResponse("test", "correlation")

	assert.Equal(t, "correlation", testDialer.connection.channel.lastCorrelationId)
}

func TestNegativeAcknowledgeMessage(t *testing.T) {
	setupDialer()
	mq, _ := NewMessageQueue()
	message := interfaces.Message{Id: 1}

	mq.Nack(message)

	assert.Equal(t, uint64(1), testDialer.connection.channel.nackedMessageId)
	assert.True(t, testDialer.connection.channel.nackRequeue)
}
