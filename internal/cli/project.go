package cli

import (
	"fmt"

	"github.com/CnTeng/table"
	"github.com/CnTeng/todoist-api-go/sync"
	"github.com/CnTeng/todoist-cli/internal/utils"
	"github.com/fatih/color"
)

func (c *cli) PrintProjects(ps []*sync.Project) {
	tbl := table.NewTable()
	tbl.AddHeader("ID", "Name", "Color")
	tbl.SetHeaderStyle(headerStyle)
	tbl.SetColStyle(1, &table.CellStyle{WrapText: utils.BoolPtr(true)})

	for _, p := range ps {
		row := table.Row{}

		row = append(row, p.ID)
		row = append(row, c.projectName(p))
		row = append(row, color.BgRGB(p.Color.RGB()).Sprint(p.Color))

		tbl.AddRow(row)
	}

	fmt.Print(tbl.Render())
}

func (c *cli) projectName(p *sync.Project) *table.Cell {
	icon := c.icons.none
	if p.InboxProject {
		icon = c.icons.inbox
	} else if p.IsFavorite {
		icon = c.icons.favorite
	}
	return &table.Cell{
		Content: p.Name,
		PrefixFunc: func(isFirst, isLast bool) string {
			if isFirst {
				return icon
			}
			return c.icons.none
		},
	}
}
