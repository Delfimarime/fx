package messaging

import (
	"fmt"
	"github.com/delfimarime/fx/config"
	"strings"
)

// Security mechanism constants
const (
	TokenAuthenticationType  = "TOKEN"
	TLSAuthenticationType    = "SSL/TLS"
	ClientAuthenticationType = "SASL/PLAIN"
	NoneAuthenticationType   = "NONE"
)

func newSecurityMechanismError(config config.Channel, channelType string, opts ...string) error {
	return fmt.Errorf(`%s supports the following security mechanisms: [%s]. Provided: "%s"`,
		channelType, strings.Join(opts, ", "), config.ChannelSecurity.Mechanism)
}
