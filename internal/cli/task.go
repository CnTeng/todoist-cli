package cli

import (
	"fmt"
	"strings"

	"github.com/CnTeng/table"
	"github.com/CnTeng/todoist-cli/internal/model"
	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/text"
)

func (c *cli) PrintTasks(ts []*model.Task) {
	tbl := table.NewTable()
	tbl.AddHeader("ID", "  ", "Project", "Name", "Description", "Labels", "Due")
	tbl.SetHeaderStyle(&table.CellStyle{CellAttrs: text.Colors{text.FgGreen, text.Underline}})

	for _, t := range ts {
		row := []any{}

		row = append(row, t.ID)

		priorityColor := text.FgWhite

		if t.CompletedAt != nil {
			row = append(row, priorityColor.Sprint(c.icons.done))
		} else {
			row = append(row, priorityColor.Sprint(c.icons.undone))
		}

		row = append(row, t.Project)

		labels := make([]string, len(t.Labels))
		if t.Labels != nil {
			for i, l := range t.Labels {
				labels[i] = color.BgRGB(l.Color.RGB()).Sprint(l.Name)
			}
		}

		due := ""
		if t.Due != nil {
			due = t.Due.String
		}

		row = append(
			row,
			t.Content,
			t.Description,
			strings.Join(labels, ", "),
			due,
		)

		tbl.AddRow(row...)
	}

	fmt.Print(tbl.Render())
}
