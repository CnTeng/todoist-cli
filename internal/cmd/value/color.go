package value

import "github.com/CnTeng/todoist-api-go/sync"

type ColorPtr struct {
	value **sync.Color
}

func NewColorPtr(value **sync.Color) *ColorPtr {
	return &ColorPtr{value: value}
}

func (c *ColorPtr) Set(val string) error {
	color, err := sync.ParseColor(val)
	if err != nil {
		return err
	}
	*c.value = &color
	return nil
}

func (c *ColorPtr) Type() string {
	return "string"
}

func (c *ColorPtr) String() string {
	if *c.value == nil {
		return ""
	}
	return string(**c.value)
}
