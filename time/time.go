package time

import (
	"encoding/json"
	"time"
)

const (
	ISOTimeFormat     = "15:04:05"
	ISODateFormat     = "2006-01-02"
	ISODatetimeFormat = "2006-01-02T15:04:05"
)

type ISOTime time.Time
type ISODate time.Time
type ISODatetime time.Time

func (it ISOTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(it).Format(ISOTimeFormat))
}

func (it *ISOTime) UnmarshalJSON(data []byte) error {
	var timeStr string
	if err := json.Unmarshal(data, &timeStr); err != nil {
		return err
	}
	t, err := time.Parse(ISOTimeFormat, timeStr)
	if err != nil {
		return err
	}
	*it = ISOTime(t)
	return nil
}

func (id ISODate) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(id).Format(ISODateFormat))
}

func (id *ISODate) UnmarshalJSON(data []byte) error {
	var dateStr string
	if err := json.Unmarshal(data, &dateStr); err != nil {
		return err
	}
	d, err := time.Parse(ISODateFormat, dateStr)
	if err != nil {
		return err
	}
	*id = ISODate(d)
	return nil
}

func (idt ISODatetime) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(idt).Format(ISODatetimeFormat))
}

func (idt *ISODatetime) UnmarshalJSON(data []byte) error {
	var datetimeStr string
	if err := json.Unmarshal(data, &datetimeStr); err != nil {
		return err
	}
	dt, err := time.Parse(ISODatetimeFormat, datetimeStr)
	if err != nil {
		return err
	}
	*idt = ISODatetime(dt)
	return nil
}

func (it ISOTime) MarshalYAML() (interface{}, error) {
	return time.Time(it).Format(ISOTimeFormat), nil
}

func (it *ISOTime) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var timeStr string
	if err := unmarshal(&timeStr); err != nil {
		return err
	}
	t, err := time.Parse(ISOTimeFormat, timeStr)
	if err != nil {
		return err
	}
	*it = ISOTime(t)
	return nil
}

func (id ISODate) MarshalYAML() (interface{}, error) {
	return time.Time(id).Format(ISODateFormat), nil
}

func (id *ISODate) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var dateStr string
	if err := unmarshal(&dateStr); err != nil {
		return err
	}
	d, err := time.Parse(ISODateFormat, dateStr)
	if err != nil {
		return err
	}
	*id = ISODate(d)
	return nil
}

func (idt ISODatetime) MarshalYAML() (interface{}, error) {
	return time.Time(idt).Format(ISODatetimeFormat), nil
}

func (idt *ISODatetime) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var datetimeStr string
	if err := unmarshal(&datetimeStr); err != nil {
		return err
	}
	dt, err := time.Parse(ISODatetimeFormat, datetimeStr)
	if err != nil {
		return err
	}
	*idt = ISODatetime(dt)
	return nil
}

func (it *ISODatetime) EqualTo(other time.Time) bool {
	src := time.Time(*it) // converting ISODatetime to time.Time
	return src.Equal(other)
}
