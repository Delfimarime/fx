package config

type Config struct {
	Name    string                    `json:"name" yaml:"name" xml:"name"`
	Server  Server                    `json:"server" yaml:"server" xml:"server"`
	Gorm    map[string]GormDatasource `json:"gorm,omitempty" yaml:"gorm,omitempty"`
	Channel map[string]Channel        `json:"channel,omitempty" yaml:"channel,omitempty"`
}

type Server struct {
	Type        string   `json:"type" yaml:"type" xml:"type"`
	Port        int      `json:"port" yaml:"port" xml:"port"`
	Accept      []string `json:"accept" yaml:"accept" xml:"accept"`
	Enabled     bool     `json:"is_enabled" yaml:"is_enabled" xml:"is_enabled"`
	ContentType []string `json:"content-type" yaml:"content-type" xml:"content-type"`
}

type Channel struct {
	Kafka
	Pulsar
	Rabbitmq
	ChannelAuthentication
	Topic      string     `json:"topic,omitempty" yaml:"topic,omitempty"`
	TLSOptions TLSOptions `json:"tls,omitempty" yaml:"tls,omitempty"`
	Type       string     `json:"type,omitempty" yaml:"type,omitempty"`
}

type Rabbitmq struct {
	Host       string `json:"host,omitempty" yaml:"host,omitempty"`
	Exchange   string `json:"exchange,omitempty" yaml:"exchange,omitempty"`
	RoutingKey string `json:"routing_key,omitempty" yaml:"routing_key,omitempty"`
}

type Kafka struct {
	Host string `json:"host,omitempty" yaml:"host,omitempty"`
}

type Pulsar struct {
	Host           string `json:"host,omitempty" yaml:"host,omitempty"`
	Token          string `json:"token,omitempty" yaml:"token,omitempty"`
	TrustCertsFile File   `json:"trust_certs,omitempty" yaml:"trust_certs,omitempty"`
}

type ChannelAuthentication struct {
	Username          string `json:"username,omitempty" yaml:"username,omitempty"`
	Password          string `json:"password,omitempty" yaml:"password,omitempty"`
	ClientKey         File   `json:"client_key,omitempty" yaml:"client_key,omitempty"`
	ClientCertificate File   `json:"client_certificate,omitempty" yaml:"certificate,omitempty"`
	SecurityMechanism string `json:"security_mechanism,omitempty" yaml:"security_mechanism,omitempty"`
}

type TLSOptions struct {
	RootCA                   File `json:"certificate,omitempty" yaml:"certificate,omitempty"`
	SkipVerification         bool `json:"skip_verification,omitempty" yaml:"skip_verification,omitempty"`
	SkipHostnameVerification bool `json:"skip_hostname_verification,omitempty" yaml:"skip_hostname_verification,omitempty"`
}

type GormDatasource struct {
	Type           string                       `json:"type,omitempty" yaml:"type,omitempty"`
	URL            string                       `json:"url,omitempty" yaml:"url,omitempty"`
	Host           string                       `json:"host,omitempty" yaml:"host,omitempty"`
	Database       string                       `json:"database,omitempty" yaml:"database,omitempty"`
	Authentication GormDatasourceAuthentication `json:"authentication,omitempty" yaml:"authentication,omitempty"`
}

type GormDatasourceAuthentication struct {
	Username string `json:"username,omitempty" yaml:"username,omitempty"`
	Password string `json:"password,omitempty" yaml:"password,omitempty"`
}

type File struct {
	URI string `json:"uri,omitempty" yaml:"uri,omitempty"`
}
