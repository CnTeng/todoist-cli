package view

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/CnTeng/todoist-api-go/sync"
	"github.com/CnTeng/todoist-cli/internal/model"
)

type labelReorderView struct {
	icons  *Icons
	labels []*model.Label
	params *sync.LabelReorderArgs
}

func NewReorderLabelView(labels []*model.Label, icons *Icons, param *sync.LabelReorderArgs) InteractiveView {
	return &labelReorderView{
		icons:  icons,
		labels: labels,
		params: param,
	}
}

func (v *labelReorderView) Render() string {
	b := &strings.Builder{}

	for _, l := range v.labels {
		if l.IsShared {
			continue
		}
		fmt.Fprintf(b, "%s %s\n", l.ID, l.Name)
	}

	return b.String()
}

func (v *labelReorderView) Interact() error {
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
		if len(fields) < 2 {
			return fmt.Errorf("invalid line %d: %s", index+1, line)
		}

		id := fields[0]
		v.params.IDOrderMapping[string(id)] = index
		index++
	}

	return nil
}
