package mq

import (
	"fmt"
	"github.com/c3sr/mq/interfaces"
	"log"
	"os"
)

var dialer interfaces.QueueDialer

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func SetDialer(d interfaces.QueueDialer) {
	dialer = d
}

func NewMessageQueue() (mq interfaces.MessageQueue, err error) {
	if dialer == nil {
		err = fmt.Errorf("A dialer must be provided using SetDialer()")
		return
	}

	url := makeMqUrlFromEnvironment()

	var ch interfaces.QueueChannel
	var conn interfaces.QueueConnection

	conn, err = dialer.Dial(url)
	failOnError(err, "Failed to connect to message queue")
	ch, err = conn.Channel()
	failOnError(err, "Failed to open a channel")

	mq = &messageQueue{
		channel:    ch,
		connection: conn,
	}

	return
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
