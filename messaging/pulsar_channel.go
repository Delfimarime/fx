package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/apache/pulsar/pulsar-client-go/pulsar"
	"github.com/delfimarime/fx/config"
	"go.uber.org/zap"
)

type PulsarProducer interface {
	Close() error
	SendAndGetMsgID(context.Context, pulsar.ProducerMessage) (pulsar.MessageID, error)
}

type PulsarChannel struct {
	producer PulsarProducer
}

func (instance *PulsarChannel) GetType() string {
	return "pulsar"
}

func (instance *PulsarChannel) Close() error {
	err := instance.producer.Close()
	if err != nil {
		return err
	}
	return nil
}

func (instance *PulsarChannel) Accept(event Event) error {
	bytes, err := json.Marshal(event)
	if err != nil {
		zap.L().Error("cannot marshal event",
			zap.String("eventId", event.Id),
			zap.String("eventType", event.Type),
			zap.String("component", "pulsar_channel"),
			zap.Error(err),
		)
		return err
	}
	messageId, err := instance.producer.SendAndGetMsgID(context.TODO(), pulsar.ProducerMessage{
		Payload: bytes,
	})
	if err != nil {
		zap.L().Error("cannot send message",
			zap.String("eventId", event.Id),
			zap.String("eventType", event.Type),
			zap.String("component", "pulsar_channel"),
			zap.ByteString("event", bytes),
			zap.Error(err),
		)
		return err
	}
	zap.L().Info("message sent",
		zap.String("eventId", event.Id),
		zap.String("eventType", event.Type),
		zap.String("component", "pulsar_channel"),
		zap.ByteString("pulsar_id", messageId.Serialize()),
	)
	return nil
}

func NewPulsarChannel(config config.Channel) (Channel, error) {
	opts := pulsar.ClientOptions{
		URL: config.Host,
	}
	err := configureTLSForPulsar(&opts, config)
	if err != nil {
		return nil, err
	}
	switch config.ChannelSecurity.Mechanism {
	case TokenAuthenticationType:
		opts.Authentication = pulsar.NewAuthenticationToken(config.Pulsar.Token)
	case TLSAuthenticationType:
		// Handled inside the configureTLSForPulsar function
	default:
		return nil, newSecurityMechanismError(config, TokenAuthenticationType, TLSAuthenticationType)
	}
	client, err := pulsar.NewClient(opts)
	if err != nil {
		return nil, err
	}
	producer, err := client.CreateProducer(pulsar.ProducerOptions{Topic: config.Topic})
	if err != nil {
		return nil, err
	}
	return &PulsarChannel{producer: producer}, nil
}

func configureTLSForPulsar(opts *pulsar.ClientOptions, config config.Channel) error {
	if config.TLSOptions.RootCA.URI != "" || config.ChannelSecurity.Mechanism == TLSAuthenticationType {
		opts.TLSTrustCertsFilePath = config.TLSOptions.RootCA.URI
		opts.URL = fmt.Sprintf("pulsar://%s", opts.URL)
	} else {
		opts.URL = fmt.Sprintf("pulsar+ssl://%s", opts.URL)
	}
	if config.ChannelSecurity.Mechanism == TLSAuthenticationType {
		opts.TLSAllowInsecureConnection = config.TLSOptions.SkipVerification
		opts.Authentication = pulsar.NewAuthenticationTLS(config.ChannelSecurity.ClientCertificate.URI, config.ChannelSecurity.ClientKey.URI)
	}
	return nil
}
