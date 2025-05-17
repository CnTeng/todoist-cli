package view

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/CnTeng/todoist-api-go/sync"
)

type filterReorderView struct {
	icons   *Icons
	filters []*sync.Filter
	params  *sync.FilterReorderArgs
}

func NewFilterReorderView(filters []*sync.Filter, icons *Icons, param *sync.FilterReorderArgs) InteractiveView {
	return &filterReorderView{
		icons:   icons,
		filters: filters,
		params:  param,
	}
}

func (v *filterReorderView) Render() string {
	b := &strings.Builder{}

	for _, f := range v.filters {
		fmt.Fprintf(b, "%s %s %s\n", f.ID, f.Name, f.Query)
	}

	return b.String()
}

func (v *filterReorderView) Interact() error {
	data, err := tempFile(v.Render())
	if err != nil {
		return err
	}

	index := 0
	for line := range bytes.Lines(data) {
		if len(line) == 0 && bytes.HasPrefix(line, []byte("#")) {
			continue
		}

		fields := bytes.Fields(line)
		if len(fields) < 3 {
			return fmt.Errorf("invalid line %d: %s", index+1, line)
		}

		id := fields[0]
		v.params.IDOrderMapping[string(id)] = index
		index++
	}

	return nil
}
