package value

import (
	"time"

	"github.com/spf13/pflag"
)

type timeValue struct {
	value  time.Time
	layout string
	action func(v time.Time) error
}

func NewTimeValue(layout string, action func(v time.Time) error) pflag.Value {
	return &timeValue{
		layout: layout,
		action: action,
	}
}

func (t *timeValue) Set(val string) error {
	v, err := time.Parse(t.layout, val)
	if err != nil {
		return err
	}
	t.value = v
	return t.action(t.value)
}

func (t *timeValue) Type() string {
	return "time"
}

func (t *timeValue) String() string {
	if t.value.IsZero() {
		return ""
	}
	return t.value.Format(time.DateOnly)
}
