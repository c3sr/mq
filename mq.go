package mq

import (
	"fmt"
	"log"
	"mq/interfaces"
	"os"
)

var connection interfaces.QueueConnection
var qChannel interfaces.QueueChannel
var dialer interfaces.QueueDialer

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func SetDialer(d interfaces.QueueDialer) {
	dialer = d
}

func NewMessageQueue() interfaces.MessageQueue {
	url := makeMqUrlFromEnvironment()

	if connection == nil {
		connection, err := dialer.Dial(url)
		failOnError(err, "Failed to connect to message queue")
		qChannel, err = connection.Channel()
		failOnError(err, "Failed to open a channel")
	}

	return messageQueue{
		channel: qChannel,
	}
}

func makeMqUrlFromEnvironment() string {
	return fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		os.Getenv("MQ_USER"),
		os.Getenv("MQ_PASSWORD"),
		os.Getenv("MQ_HOST"),
		os.Getenv("MQ_PORT"),
	)
}
