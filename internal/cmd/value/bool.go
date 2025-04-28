package value

import "strconv"

type BoolPtr struct {
	value **bool
}

func NewBoolPtr(value **bool) *BoolPtr {
	return &BoolPtr{value: value}
}

func (b *BoolPtr) Set(val string) error {
	v, err := strconv.ParseBool(val)
	if err != nil {
		return err
	}
	*b.value = &v
	return nil
}

func (b *BoolPtr) Type() string {
	return "bool"
}

func (b *BoolPtr) String() string {
	if *b.value == nil {
		return ""
	}
	return strconv.FormatBool(**b.value)
}
