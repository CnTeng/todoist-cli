package view

import (
	"github.com/CnTeng/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

type View interface {
	Render() string
}

var headerStyle = &table.CellStyle{
	CellAttrs: text.Colors{text.FgGreen, text.Underline},
}
