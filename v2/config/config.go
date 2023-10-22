package config

type Config struct {
	Data         Data               `json:"data" yaml:"data" xml:"data"`
	Server       Server             `json:"server" yaml:"server" xml:"server"`
	Integrations map[string]Channel `json:"integrations" yaml:"integrations" xml:"integrations"`
}

type Data struct {
	Gorm *GormDb `json:"gorm,omitempty" yaml:"gorm,omitempty"`
}

type Server struct {
	Mode        string   `json:"mode" yaml:"mode" xml:"mode"`
	Type        string   `json:"type" yaml:"type" xml:"type"`
	Port        int      `json:"port" yaml:"port" xml:"port"`
	Accept      []string `json:"accept" yaml:"accept" xml:"accept"`
	ContentType []string `json:"content-type" yaml:"content-type" xml:"content-type"`
}

type Channel struct {
	Pulsar          `yaml:",inline"`
	Rabbitmq        `yaml:",inline"`
	TLSOptions      TLSOptions      `json:"tls,omitempty" yaml:"tls,omitempty"`
	Host            string          `json:"host,omitempty" yaml:"host,omitempty"`
	Type            string          `json:"type,omitempty" yaml:"type,omitempty"`
	Topic           string          `json:"topic,omitempty" yaml:"topic,omitempty"`
	ChannelSecurity ChannelSecurity `json:"authentication,omitempty" yaml:"authentication,omitempty"`
}

type Rabbitmq struct {
	Exchange   string `json:"exchange,omitempty" yaml:"exchange,omitempty"`
	RoutingKey string `json:"routing_key,omitempty" yaml:"routing_key,omitempty"`
}

type Pulsar struct {
	Token          string `json:"token,omitempty" yaml:"token,omitempty"`
	TrustCertsFile File   `json:"trust_certs,omitempty" yaml:"trust_certs,omitempty"`
}

type ChannelSecurity struct {
	BasicAuthentication `yaml:",inline"`
	ClientKey           File   `json:"client_key,omitempty" yaml:"client_key,omitempty"`
	ClientCertificate   File   `json:"client_certificate,omitempty" yaml:"certificate,omitempty"`
	Mechanism           string `json:"security_mechanism,omitempty" yaml:"security_mechanism,omitempty"`
}

type TLSOptions struct {
	RootCA                   File `json:"certificate,omitempty" yaml:"certificate,omitempty"`
	SkipVerification         bool `json:"skip_verification,omitempty" yaml:"skip_verification,omitempty"`
	SkipHostnameVerification bool `json:"skip_hostname_verification,omitempty" yaml:"skip_hostname_verification,omitempty"`
}

type GormDb struct {
	Type           string              `json:"type,omitempty" yaml:"type,omitempty"`
	URL            string              `json:"url,omitempty" yaml:"url,omitempty"`
	Host           string              `json:"host,omitempty" yaml:"host,omitempty"`
	Database       string              `json:"database,omitempty" yaml:"database,omitempty"`
	Authentication BasicAuthentication `json:"authentication,omitempty" yaml:"authentication,omitempty"`
}

type BasicAuthentication struct {
	Username string `json:"username,omitempty" yaml:"username,omitempty"`
	Password string `json:"password,omitempty" yaml:"password,omitempty"`
}

type File struct {
	URI string `json:"uri,omitempty" yaml:"uri,omitempty"`
}
