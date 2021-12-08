# MLModelScope Message Queue
This library provides an abstraction layer for a lower level message queue implementation. The goal is to provide a
simple interface for interacting with a message queue in MLModelScope.

## Requirements

This library requires go version >= 1.13

## Supported Message Queues

Support is currently provided for:

* RabbitMQ (AMQP 0.9.1)

## Configuration

When a message queue is created using `mq.NewMessageQueue`, configuration is read from environment variables. The
following variables must be set for a message queue to be configured correctly:

* `MQ_HOST` _the hostname of the message queue server_
* `MQ_PORT` _the port number the message queue server listens on_
* `MQ_USER` _a username that can login to the message queue server_
* `MQ_PASSWORD` _the password that corresponds to the username_

## Integration tests

A small suite of integration tests are included for the supported message queue implementation(s). These tests require
a functional message queue server to connect to. Configuration for connecting to the message queue must be provided
in the environment variables listed above.

The script at `/scripts/run-integration-tests.sh` will start a RabbitMQ server in a docker container and run the
integration tests against it. This requires that you have Docker and at least go version 1.13 installed.

## Code Coverage

The script at `/scripts/generate-coverage-report.sh` will run the Unit and Integration tests with coverage, and produce
a report `coverage.html` from the merged results.