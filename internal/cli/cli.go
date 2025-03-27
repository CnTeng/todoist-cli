package cli

import (
	"github.com/CnTeng/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

type cli struct {
	icons *icons
}

func NewCLI(iconType iconType) *cli {
	return &cli{icons: newIcons(iconType).withColor()}
}

var headerStyle = &table.CellStyle{
	CellAttrs: text.Colors{text.FgGreen, text.Underline},
}
