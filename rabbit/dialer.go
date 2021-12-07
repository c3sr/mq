package rabbit

import (
	"errors"
	"fmt"
	"github.com/c3sr/mq/interfaces"
	amqp "github.com/rabbitmq/amqp091-go"
	"os"
)

type DialerConfigurationError struct {
	Field string
	Err   error
}

func (e *DialerConfigurationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Err, e.Field)
}

type rabbitDialer struct {
	url string
}

func (d *rabbitDialer) Dial() (interfaces.QueueConnection, error) {
	conn, err := amqp.Dial(d.url)

	return &rabbitConnection{conn}, err
}

func (d *rabbitDialer) URL() string {
	return d.url
}

func NewRabbitDialer() (d interfaces.QueueDialer, err error) {
	url, err := makeMqUrlFromEnvironment()
	if err == nil {
		d = &rabbitDialer{
			url: url,
		}
	}

	return
}

func makeMqUrlFromEnvironment() (url string, err error) {
	if os.Getenv("MQ_HOST") == "" {
		err = &DialerConfigurationError{
			Field: "MQ_HOST",
			Err:   errors.New("missing RabbitMQ dialer environment variable"),
		}

		return
	}

	if os.Getenv("MQ_PORT") == "" {
		err = &DialerConfigurationError{
			Field: "MQ_PORT",
			Err:   errors.New("missing RabbitMQ dialer environment variable"),
		}

		return
	}

	if os.Getenv("MQ_USER") == "" {
		err = &DialerConfigurationError{
			Field: "MQ_USER",
			Err:   errors.New("missing RabbitMQ dialer environment variable"),
		}

		return
	}

	if os.Getenv("MQ_PASSWORD") == "" {
		err = &DialerConfigurationError{
			Field: "MQ_PASSWORD",
			Err:   errors.New("missing RabbitMQ dialer environment variable"),
		}

		return
	}

	url = fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		os.Getenv("MQ_USER"),
		os.Getenv("MQ_PASSWORD"),
		os.Getenv("MQ_HOST"),
		os.Getenv("MQ_PORT"),
	)

	return
}
