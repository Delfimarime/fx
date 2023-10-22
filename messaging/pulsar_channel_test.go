package messaging

import (
	"context"
	"errors"
	"github.com/apache/pulsar/pulsar-client-go/pulsar"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockPulsarProducer struct {
	mock.Mock
}

func (m *MockPulsarProducer) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockPulsarProducer) SendAndGetMsgID(ctx context.Context, msg pulsar.ProducerMessage) (pulsar.MessageID, error) {
	args := m.Called(ctx, msg)
	return args.Get(0).(pulsar.MessageID), args.Error(1)
}

func TestPulsarChannel_Accept(t *testing.T) {
	var mockProducer *MockPulsarProducer
	var pulsarChannel *PulsarChannel
	setup := func() {
		mockProducer = new(MockPulsarProducer)
		pulsarChannel = &PulsarChannel{
			producer: mockProducer,
		}
		defer mockProducer.AssertExpectations(t)
	}
	t.Run("accept with successful message send", func(t *testing.T) {
		setup()
		event := Message{Id: "123", Type: "test"}
		mockProducer.On("SendAndGetMsgID", mock.Anything, mock.Anything).Return(&mockMessageID{}, nil)

		err := pulsarChannel.Accept(event)
		assert.Nil(t, err)
		mockProducer.AssertExpectations(t)
	})
	t.Run("accept with send message error", func(t *testing.T) {
		setup()
		event := Message{Id: "123", Type: "test"}
		mockProducer.On("SendAndGetMsgID", mock.Anything, mock.Anything).Return(mockMessageID{}, errors.New("send error"))
		err := pulsarChannel.Accept(event)
		assert.Equal(t, "send error", err.Error())
		mockProducer.AssertExpectations(t)
	})
}
func TestPulsarChannel_Close(t *testing.T) {
	var mockProducer *MockPulsarProducer
	var pulsarChannel *PulsarChannel
	setup := func() {
		mockProducer = new(MockPulsarProducer)
		pulsarChannel = &PulsarChannel{
			producer: mockProducer,
		}
		defer mockProducer.AssertExpectations(t)
	}
	t.Run("close successfully", func(t *testing.T) {
		setup()
		mockProducer.On("Close").Return(nil)
		err := pulsarChannel.Close()
		assert.Nil(t, err)
		mockProducer.AssertExpectations(t)
	})

	t.Run("close with error", func(t *testing.T) {
		setup()
		mockProducer.On("Close").Return(errors.New("close error"))

		err := pulsarChannel.Close()
		assert.Equal(t, "close error", err.Error())
		mockProducer.AssertExpectations(t)
	})
}

type mockMessageID struct{}

func (m mockMessageID) Serialize() []byte {
	return []byte("mocked_id")
}
