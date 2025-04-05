package cli

import (
	"github.com/CnTeng/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

type Cli struct {
	icons *icons
}

func NewCLI(iconType iconType) *Cli {
	return &Cli{icons: newIcons(iconType).withColor()}
}

var headerStyle = &table.CellStyle{
	CellAttrs: text.Colors{text.FgGreen, text.Underline},
}
