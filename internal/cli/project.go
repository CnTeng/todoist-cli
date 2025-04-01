package cli

import (
	"fmt"

	"github.com/CnTeng/table"
	"github.com/CnTeng/todoist-api-go/sync"
	"github.com/fatih/color"
)

func (c *cli) PrintProjects(ps []*sync.Project) {
	tbl := table.NewTable()
	tbl.AddHeader("ID", "  ", "Name", "Color")
	tbl.SetHeaderStyle(headerStyle)

	for _, p := range ps {
		row := []any{}

		row = append(row, p.ID)
		icon := " "
		if p.InboxProject {
			icon = c.icons.inbox
		} else if p.IsFavorite {
			icon = c.icons.favorite
		}
		row = append(row, icon)
		row = append(row, p.Name)
		row = append(row, color.BgRGB(p.Color.RGB()).Sprint(p.Color))

		tbl.AddRow(row...)
	}

	fmt.Print(tbl.Render())
}
