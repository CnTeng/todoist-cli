package view

import (
	"github.com/CnTeng/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

type View interface {
	Render() string
}

type ViewConfig struct {
	*TaskViewConfig
}

var DefaultViewConfig = &ViewConfig{
	TaskViewConfig: &TaskViewConfig{
		Completed:   false,
		Description: false,
		Tree:        false,
	},
}

var headerStyle = &table.CellStyle{
	CellAttrs: text.Colors{text.FgGreen, text.Underline},
}

func boolPtr(v bool) *bool {
	return &v
}
