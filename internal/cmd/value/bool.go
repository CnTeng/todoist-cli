package value

import (
	"strconv"

	"github.com/spf13/pflag"
)

type boolValue struct {
	value  bool
	action func(v bool) error
}

func NewBoolPtr(action func(v bool) error) pflag.Value {
	return &boolValue{action: action}
}

func (b *boolValue) Set(val string) error {
	v, err := strconv.ParseBool(val)
	if err != nil {
		return err
	}
	b.value = v
	return b.action(b.value)
}

func (b *boolValue) Type() string {
	return "bool"
}

func (b *boolValue) String() string {
	return strconv.FormatBool(b.value)
}
