package value

import (
	"github.com/CnTeng/todoist-api-go/sync"
)

type DurationPtr struct {
	value **sync.Duration
}

func NewDurationPtr(value **sync.Duration) *DurationPtr {
	return &DurationPtr{value: value}
}

func (d *DurationPtr) Set(val string) error {
	duration, err := sync.ParseDuration(val)
	if err != nil {
		return err
	}
	*d.value = duration
	return nil
}

func (d *DurationPtr) Type() string {
	return "string"
}

func (d *DurationPtr) String() string {
	if *d.value == nil {
		return ""
	}
	return (*d.value).String()
}
