package messaging

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/delfimarime/fx/config"
	"go.uber.org/zap"
	"os"
)

type KafkaProducer interface {
	Close() error
	SendMessage(msg *sarama.ProducerMessage) (partition int32, offset int64, err error)
}

type KafkaChannel struct {
	topic    string
	producer KafkaProducer
}

func (instance *KafkaChannel) GetType() string {
	return "kafka"
}

func (instance *KafkaChannel) Close() error {
	err := instance.producer.Close()
	if err != nil {
		return err
	}
	return nil
}

func (instance *KafkaChannel) Accept(event Event) error {
	binary, err := json.Marshal(event)
	if err != nil {
		zap.L().Error("cannot marshal event",
			zap.String("eventId", event.Id),
			zap.String("eventType", event.Type),
			zap.String("component", "kafka_channel"),
			zap.Error(err),
		)
		return err
	}
	partition, offset, err := instance.producer.SendMessage(&sarama.ProducerMessage{
		Topic: instance.topic,
		Value: sarama.StringEncoder(binary),
	})
	if err != nil {
		zap.L().Error("cannot send message to kafka topic",
			zap.String("eventId", event.Id),
			zap.String("eventType", event.Type),
			zap.String("component", "kafka_channel"),
			zap.ByteString("event", binary),
			zap.String("topic", instance.topic),
			zap.Error(err),
		)
		return err
	}
	zap.L().Info("message sent to kafka topic",
		zap.String("eventId", event.Id),
		zap.String("eventType", event.Type),
		zap.String("topic", instance.topic),
		zap.String("component", "kafka_channel"),
		zap.Int32("kafka_partition", partition),
		zap.Int64("kafka_offset", offset),
	)
	return nil
}

func NewKafkaChannel(config config.Channel) (Channel, error) {
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Net.TLS.Enable = true
	err := configureTLSForKafka(kafkaConfig, config)
	if err != nil {
		return nil, err
	}
	switch config.ChannelSecurity.Mechanism {
	case ClientAuthenticationType:
		kafkaConfig.Net.SASL.User = config.ChannelSecurity.Username
		kafkaConfig.Net.SASL.Password = config.ChannelSecurity.Password
	case TLSAuthenticationType:
		// Handled inside the configureTLSForKafka function
	default:
		return nil, newSecurityMechanismError(config, ClientAuthenticationType, TLSAuthenticationType)
	}
	producer, err := sarama.NewSyncProducer([]string{config.Host}, kafkaConfig)
	if err != nil {
		return nil, err
	}
	return &KafkaChannel{producer: producer, topic: config.Topic}, nil
}

func configureTLSForKafka(kafkaConfig *sarama.Config, config config.Channel) error {
	if config.ChannelSecurity.Mechanism == TLSAuthenticationType {
		cert, err := tls.LoadX509KeyPair(config.ChannelSecurity.ClientCertificate.URI, config.ChannelSecurity.ClientKey.URI)
		if err != nil {
			return err
		}
		kafkaConfig.Net.TLS.Config = &tls.Config{Certificates: []tls.Certificate{cert}}
	}

	if config.TLSOptions.RootCA.URI != "" {
		caCert, err := os.ReadFile(config.TLSOptions.RootCA.URI)
		if err != nil {
			return err
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		if kafkaConfig.Net.TLS.Config == nil {
			kafkaConfig.Net.TLS.Config = &tls.Config{}
		}
		kafkaConfig.Net.TLS.Config.RootCAs = caCertPool
	}

	if config.TLSOptions.SkipVerification {
		if kafkaConfig.Net.TLS.Config == nil {
			kafkaConfig.Net.TLS.Config = &tls.Config{}
		}
		kafkaConfig.Net.TLS.Config.InsecureSkipVerify = true
	}
	return nil
}
