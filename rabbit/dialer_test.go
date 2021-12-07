// +build !integration

package rabbit

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNewRabbitDialerReturnsErrorForMissingMqHost(t *testing.T) {
	os.Clearenv()
	os.Setenv("MQ_PORT", "1234")
	os.Setenv("MQ_USER", "testuser")
	os.Setenv("MQ_PASSWORD", "testpassword")

	_, err := NewRabbitDialer()

	configError, ok := err.(*DialerConfigurationError)
	assert.True(t, ok, "NewRabbitDialer() should return a DialerConfigurationError")
	assert.Equal(t, "MQ_HOST", configError.Field)
	assert.Equal(t, "missing RabbitMQ dialer environment variable: MQ_HOST", err.Error())
}

func TestNewRabbitDialerReturnsErrorForMissingMqPort(t *testing.T) {
	os.Clearenv()
	os.Setenv("MQ_HOST", "testhost")
	os.Setenv("MQ_USER", "testuser")
	os.Setenv("MQ_PASSWORD", "testpassword")

	_, err := NewRabbitDialer()

	configError, ok := err.(*DialerConfigurationError)
	assert.True(t, ok, "NewRabbitDialer() should return a DialerConfigurationError")
	assert.Equal(t, "MQ_PORT", configError.Field)
}

func TestNewRabbitDialerReturnsErrorForMissingMqUser(t *testing.T) {
	os.Clearenv()
	os.Setenv("MQ_HOST", "testhost")
	os.Setenv("MQ_PORT", "1234")
	os.Setenv("MQ_PASSWORD", "testpassword")

	_, err := NewRabbitDialer()

	configError, ok := err.(*DialerConfigurationError)
	assert.True(t, ok, "NewRabbitDialer() should return a DialerConfigurationError")
	assert.Equal(t, "MQ_USER", configError.Field)
}

func TestNewRabbitDialerReturnsErrorForMissingMqPassword(t *testing.T) {
	os.Clearenv()
	os.Setenv("MQ_HOST", "testhost")
	os.Setenv("MQ_PORT", "1234")
	os.Setenv("MQ_USER", "testuser")

	_, err := NewRabbitDialer()

	configError, ok := err.(*DialerConfigurationError)
	assert.True(t, ok, "NewRabbitDialer() should return a DialerConfigurationError")
	assert.Equal(t, "MQ_PASSWORD", configError.Field)
}

func TestNewRabbitDialerCreatesURLFromEnvironment(t *testing.T) {
	os.Clearenv()
	os.Setenv("MQ_HOST", "testhost")
	os.Setenv("MQ_PORT", "1234")
	os.Setenv("MQ_USER", "testuser")
	os.Setenv("MQ_PASSWORD", "testpassword")
	dialer, _ := NewRabbitDialer()

	assert.Equal(t, "amqp://testuser:testpassword@testhost:1234/", dialer.URL())
}
