package value

import "strconv"

type IntPtr struct {
	value **int
}

func NewIntPtr(value **int) *IntPtr {
	return &IntPtr{value: value}
}

func (i *IntPtr) Set(val string) error {
	parsedValue, err := strconv.Atoi(val)
	if err != nil {
		return err
	}
	*i.value = &parsedValue
	return nil
}

func (i *IntPtr) Type() string {
	return "int"
}

func (i *IntPtr) String() string {
	if *i.value == nil {
		return ""
	}
	return strconv.Itoa(**i.value)
}
