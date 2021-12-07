// Package mq provides an abstraction layer around concrete Message Queue implementation(s).
// mq is not intended to be a generic interface around Message Queues. It is opinionated in how
// the underlying Message Queue is used and its API provides a very limited subset of potential
// functionality.
package mq

import (
	"fmt"
	"github.com/c3sr/mq/interfaces"
)

var dialer interfaces.QueueDialer

// SetDialer configures the package with a dialer that provides connections to a specific Message Queue
// implementation.
func SetDialer(d interfaces.QueueDialer) {
	dialer = d
}

// NewMessageQueue constructs and returns a new object implementing the interfaces.MessageQueue interface.
// This object uses the dialer provided by calling SetDialer to connect to the message queue server.
// SetDialer *must* be called before calling NewMessageQueue, otherwise an error will be returned.
func NewMessageQueue() (mq interfaces.MessageQueue, err error) {
	if dialer == nil {
		err = fmt.Errorf("A dialer must be provided using SetDialer()")
		return
	}

	var ch interfaces.QueueChannel
	var conn interfaces.QueueConnection

	conn, err = dialer.Dial()
	if err != nil {
		return
	}
	ch, err = conn.Channel()
	if err != nil {
		return
	}

	mq = &messageQueue{
		channel:    ch,
		connection: conn,
	}

	return
}
