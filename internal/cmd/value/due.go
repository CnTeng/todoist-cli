package value

import "github.com/CnTeng/todoist-api-go/sync"

type DuePtr struct {
	value **sync.Due
}

func NewDuePtr(value **sync.Due) *DuePtr {
	return &DuePtr{value: value}
}

func (d *DuePtr) Set(val string) error {
	*d.value = &sync.Due{String: &val}
	return nil
}

func (d *DuePtr) Type() string {
	return "due"
}

func (d *DuePtr) String() string {
	if *d.value == nil || (*d.value).String == nil {
		return ""
	}
	return *(*d.value).String
}
