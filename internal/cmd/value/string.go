package value

type StringPtr struct {
	value **string
}

func NewStringPtr(value **string) *StringPtr {
	return &StringPtr{value: value}
}

func (s *StringPtr) Set(val string) error {
	*s.value = &val
	return nil
}

func (s *StringPtr) Type() string {
	return "string"
}

func (s *StringPtr) String() string {
	if *s.value == nil {
		return ""
	}
	return **s.value
}
