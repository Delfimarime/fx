package messaging

import (
	"errors"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

// MockAmpqChannel mocks the AmpqChannel interface
type MockAmpqChannel struct {
	mock.Mock
}

func (m *MockAmpqChannel) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockAmpqChannel) Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	args := m.Called(exchange, key, mandatory, immediate, msg)
	return args.Error(0)
}

func TestRabbitmqChannel_Accept(t *testing.T) {
	// Setting up our mock channel
	mockChannel := new(MockAmpqChannel)
	channel := &RabbitmqChannel{}
	beforeTest := func() {
		mockChannel = new(MockAmpqChannel)
		channel = &RabbitmqChannel{
			exchange:   "unit-test",
			routingKey: t.Name(),
			channel:    mockChannel,
		}
		defer mockChannel.AssertExpectations(t)
	}
	event := Message{
		Id:          uuid.NewString(),
		Type:        "test",
		Domain:      "unit-testing",
		Src:         "golang",
		Version:     "1.0.0",
		Properties:  []byte(`{"first_name":"jJohn","last_name":"Doe"}`),
		TriggeredBy: "",
	}

	t.Run("test successful message send", func(t *testing.T) {
		beforeTest()
		mockChannel.On("Publish", channel.exchange, channel.routingKey, false, false, mock.Anything).Return(nil)
		err := channel.Accept(event)
		assert.Nil(t, err)
	})
	t.Run("test publish error", func(t *testing.T) {
		beforeTest()
		mockChannel.On("Publish", channel.exchange, channel.routingKey, false, false, mock.Anything).Return(errors.New("publish error"))
		err := channel.Accept(event)
		assert.NotNil(t, err)
	})
}

func TestRabbitmqChannel_Close(t *testing.T) {
	var mockChannel *MockAmpqChannel
	var channel *RabbitmqChannel

	setup := func() {
		mockChannel = new(MockAmpqChannel)
		channel = &RabbitmqChannel{
			exchange:   "unit-test-exchange",
			routingKey: "unit-test-routing",
			channel:    mockChannel,
		}
	}
	t.Run("successful close", func(t *testing.T) {
		setup()
		// Mocking a successful close
		mockChannel.On("Close").Return(nil)

		err := channel.Close()
		assert.Nil(t, err, "Expected no error on successful close")
		mockChannel.AssertExpectations(t) // Assert that Close was called on the mock
	})
	t.Run("close error", func(t *testing.T) {
		setup()
		// Mocking an error on close
		mockChannel.On("Close").Return(errors.New("close error"))

		err := channel.Close()
		assert.NotNil(t, err, "Expected an error on close")
		assert.Equal(t, "close error", err.Error(), "Error message mismatch")
		mockChannel.AssertExpectations(t) // Assert that Close was called on the mock
	})
}
