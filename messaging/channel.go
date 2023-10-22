package messaging

const (
	KafkaType    = "kafka"
	PulsarType   = "pulsar"
	RabbitmqType = "rabbitmq"
)

type Channel interface {
	Close() error
	GetType() string
	Accept(event Message) error
}
