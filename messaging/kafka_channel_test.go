package messaging

import (
	"errors"
	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockKafkaProducer struct {
	mock.Mock
}

func (m *MockKafkaProducer) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockKafkaProducer) SendMessage(msg *sarama.ProducerMessage) (partition int32, offset int64, err error) {
	args := m.Called(msg)
	return int32(0), int64(1), args.Error(2)
}

func TestKafkaChannel_Accept(t *testing.T) {
	var mockProducer *MockKafkaProducer
	var kafkaChannel *KafkaChannel
	setup := func() {
		mockProducer = new(MockKafkaProducer)
		kafkaChannel = &KafkaChannel{
			topic:    "unit-test-topic",
			producer: mockProducer,
		}
		defer mockProducer.AssertExpectations(t)
	}
	t.Run("accept with successful message send", func(t *testing.T) {
		setup()
		event := Event{Id: "123", Type: "test"}
		mockProducer.On("SendMessage", mock.Anything).Return(int32(1), int64(1), nil)

		err := kafkaChannel.Accept(event)
		assert.Nil(t, err)
		mockProducer.AssertExpectations(t)
	})
	t.Run("accept with send message error", func(t *testing.T) {
		setup()
		event := Event{Id: "123", Type: "test"}
		mockProducer.On("SendMessage", mock.Anything).Return(int32(0), int64(0), errors.New("send error"))

		err := kafkaChannel.Accept(event)
		assert.Equal(t, "send error", err.Error())
		mockProducer.AssertExpectations(t)
	})
}

func TestKafkaChannel_Close(t *testing.T) {
	var mockProducer *MockKafkaProducer
	var kafkaChannel *KafkaChannel

	setup := func() {
		mockProducer = new(MockKafkaProducer)
		kafkaChannel = &KafkaChannel{
			topic:    "unit-test-topic",
			producer: mockProducer,
		}
		defer mockProducer.AssertExpectations(t)
	}
	t.Run("close successfully", func(t *testing.T) {
		setup()
		mockProducer.On("Close").Return(nil)

		err := kafkaChannel.Close()
		assert.Nil(t, err)
		mockProducer.AssertExpectations(t)
	})
	t.Run("close with error", func(t *testing.T) {
		setup()
		mockProducer.On("Close").Return(errors.New("close error"))

		err := kafkaChannel.Close()
		assert.Equal(t, "close error", err.Error())
		mockProducer.AssertExpectations(t)
	})
}
