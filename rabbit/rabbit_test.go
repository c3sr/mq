// +build !integration

package rabbit

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestDialCreatesURLFromEnvironment(t *testing.T) {
	os.Setenv("MQ_HOST", "testhost")
	os.Setenv("MQ_PORT", "1234")
	os.Setenv("MQ_USER", "testuser")
	os.Setenv("MQ_PASSWORD", "testpassword")
	dialer := NewRabbitDialer()
	dialer.Dial()

	assert.Equal(t, "amqp://testuser:testpassword@testhost:1234/", dialer.URL())
}
