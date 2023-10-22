package messaging

import (
	"encoding/json"
	"fmt"
	"github.com/delfimarime/fx/config"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type AmpqChannel interface {
	Close() error
	Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error
}

type RabbitmqChannel struct {
	exchange   string
	routingKey string
	channel    AmpqChannel
}

func (instance *RabbitmqChannel) GetType() string {
	return "rabbitmq"
}

func (instance *RabbitmqChannel) Close() error {
	err := instance.channel.Close()
	if err != nil {
		return err
	}
	return nil
}

func (instance *RabbitmqChannel) Accept(event Message) error {
	bytes, err := json.Marshal(event)
	if err != nil {
		zap.L().Error("cannot marshal event",
			zap.String("eventId", event.Id),
			zap.String("eventType", event.Type),
			zap.String("component", "rabbitmq_channel"),
			zap.Error(err),
		)
		return err
	}
	err = instance.channel.Publish(
		instance.exchange,
		instance.routingKey,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        bytes,
		},
	)
	if err != nil {
		zap.L().Error("cannot send message to RabbitMQ",
			zap.String("eventId", event.Id),
			zap.String("eventType", event.Type),
			zap.String("exchange", instance.exchange),
			zap.String("routingKey", instance.routingKey),
			zap.String("component", "rabbitmq_channel"),
			zap.ByteString("event", bytes),
			zap.Error(err),
		)
		return err
	}
	zap.L().Info("message sent to RabbitMQ",
		zap.String("eventId", event.Id),
		zap.String("eventType", event.Type),
		zap.String("exchange", instance.exchange),
		zap.String("routingKey", instance.routingKey),
		zap.String("component", "rabbitmq_channel"),
	)
	return nil
}

func NewRabbitmqTypedFactory() TypedFactory {
	return TypedFactory{
		Type:    RabbitmqType,
		Factory: NewRabbitmqChannel,
	}
}

func NewRabbitmqChannel(config config.Channel) (Channel, error) {
	var host string
	switch config.ChannelSecurity.Mechanism {
	case NoneAuthenticationType:
		host = fmt.Sprintf("amqp://%s/", config.Host)
	case ClientAuthenticationType:
		host = fmt.Sprintf("amqp://%s:%s@%s/",
			config.ChannelSecurity.Username, config.ChannelSecurity.Password, config.Host)
	default:
		return nil, newSecurityMechanismError(config, "rabbitmq", NoneAuthenticationType, ClientAuthenticationType)
	}
	conn, err := amqp.Dial(host)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	return &RabbitmqChannel{channel: ch, routingKey: config.RoutingKey, exchange: config.Exchange}, nil
}
