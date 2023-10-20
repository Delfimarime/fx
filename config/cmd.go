package config

type Terminal struct {
	Config
	Logging Logging `json:"logging,omitempty" yaml:"logging,omitempty" xml:"logging,omitempty"`
	Service Service `json:"service,omitempty" yaml:"service,omitempty" xml:"service,omitempty"`
}

type Logging struct {
	Level string `json:"level,omitempty" yaml:"name,omitempty" xml:"name,omitempty"`
}

type Service struct {
	Id     string `json:"id,omitempty" yaml:"id,omitempty" xml:"id,omitempty"`
	Name   string `json:"name,omitempty" yaml:"name,omitempty" xml:"name,omitempty"`
	Domain string `json:"domain,omitempty" yaml:"domain,omitempty" xml:"domain,omitempty"`
}
