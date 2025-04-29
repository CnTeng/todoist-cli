package value

import (
	"time"

	"github.com/CnTeng/todoist-api-go/sync"
)

type DeadlinePtr struct {
	value **sync.Deadline
}

func NewDeadlinePtr(value **sync.Deadline) *DeadlinePtr {
	return &DeadlinePtr{value: value}
}

func (d *DeadlinePtr) Set(val string) error {
	time, err := time.Parse(time.DateOnly, val)
	if err != nil {
		return err
	}
	*d.value = &sync.Deadline{Date: &time}
	return nil
}

func (d *DeadlinePtr) Type() string {
	return "string"
}

func (d *DeadlinePtr) String() string {
	if *d.value == nil || (*d.value).Date == nil {
		return ""
	}
	return (*d.value).Date.Format(time.DateOnly)
}
