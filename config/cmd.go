package config

type Terminal struct {
	Config
	Logging `json:"logging,omitempty" yaml:"logging,omitempty" xml:"logging,omitempty"`
}

type Logging struct {
	Level string `json:"level,omitempty" yaml:"name,omitempty" xml:"name,omitempty"`
}
