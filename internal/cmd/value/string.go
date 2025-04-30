package value

import "github.com/spf13/pflag"

type stringValue struct {
	value  string
	action func(v string) error
}

func NewStringPtr(action func(v string) error) pflag.Value {
	return &stringValue{action: action}
}

func (s *stringValue) Set(val string) error {
	s.value = val
	return s.action(s.value)
}

func (s *stringValue) Type() string {
	return "string"
}

func (s *stringValue) String() string {
	return s.value
}
