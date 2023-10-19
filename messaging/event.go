package messaging

import (
	"encoding/json"
	"github.com/delfimarime/fx/time"
)

type Event struct {
	Id          string           `json:"id,omitempty" yaml:"id,omitempty"`
	Type        string           `json:"type,omitempty" yaml:"type,omitempty"`
	Domain      string           `json:"domain,omitempty" yaml:"domain,omitempty"`
	Src         string           `json:"source,omitempty" yaml:"source,omitempty"`
	Version     string           `json:"version,omitempty" yaml:"version,omitempty"`
	Properties  json.RawMessage  `json:"properties,omitempty" yaml:"properties,omitempty"`
	PublishedAt time.ISODatetime `json:"published_at,omitempty" yaml:"published_at,omitempty"`
	TriggeredBy string           `json:"triggered_by,omitempty" yaml:"triggered_by,omitempty"`
}
