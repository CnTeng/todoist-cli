package value

import (
	"strconv"

	"github.com/spf13/pflag"
)

type intValue struct {
	value  int
	action func(v int) error
}

func NewIntValue(action func(v int) error) pflag.Value {
	return &intValue{action: action}
}

func (i *intValue) Set(val string) error {
	v, err := strconv.Atoi(val)
	if err != nil {
		return err
	}
	i.value = v
	return i.action(i.value)
}

func (i *intValue) Type() string {
	return "int"
}

func (i *intValue) String() string {
	return strconv.Itoa(i.value)
}
