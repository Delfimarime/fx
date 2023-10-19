package messaging

import (
	"fmt"
	"github.com/delfimarime/fx/config"
	"strings"
)

const (
	module = "channel"
)

// Security mechanism constants
const (
	TokenAuthenticationType  = "TOKEN"
	TLSAuthenticationType    = "SSL/TLS"
	ClientAuthenticationType = "SASL/PLAIN"
	NoneAuthenticationType   = "NONE"
)

func GetChannelFrom(config config.Channel) func() (Channel, error) {
	return func() (Channel, error) {
		switch config.Type {
		case KafkaType:
			return NewKafkaChannel(config)
		case RabbitmqType:
			return NewRabbitmqChannel(config)
		case PulsarType:
			return NewPulsarChannel(config)
		default:
			return nil, fmt.Errorf(`unsupported type "%s", supported types: [%s, %s, %s]`,
				config.Type, KafkaType, PulsarType, RabbitmqType)
		}
	}
}

func newSecurityMechanismError(config config.Channel, channelType string, opts ...string) error {
	return fmt.Errorf(`%s supports the following security mechanisms: [%s]. Provided: "%s"`,
		channelType, strings.Join(opts, ", "), config.SecurityMechanism)
}
