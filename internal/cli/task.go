package cli

import (
	"fmt"
	"slices"
	"strings"

	"github.com/CnTeng/table"
	"github.com/CnTeng/todoist-cli/internal/model"
	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/text"
)

func (c *cli) taskPriorityColor(t *model.Task) text.Color {
	switch t.Priority {
	case 1:
		return text.FgWhite
	case 2:
		return text.FgBlue
	case 3:
		return text.FgYellow
	case 4:
		return text.FgRed
	default:
		return text.FgWhite
	}
}

func (c *cli) taskStatusIcon(t *model.Task) string {
	if t.CompletedAt != nil {
		return c.icons.done
	}
	return c.icons.undone
}

func (c *cli) renderTask(t *model.Task, depth []bool) []table.Row {
	project := color.RGB(t.Project.Color.RGB()).Sprint(t.Project.Name)

	labels := []string{}
	for _, l := range t.Labels {
		labels = append(labels, color.BgRGB(l.Color.RGB()).Sprint(l.Name))
	}

	due := ""
	if t.Due != nil {
		due = t.Due.String
	}

	deadline := ""
	if t.Deadline != nil {
		deadline = t.Deadline.Date.String()
	}

	depth = slices.Clone(depth)
	content := table.Cell{
		Content: t.Content,
		PrefixFunc: func(isFirst, isLast bool) string {
			b := &strings.Builder{}
			lastIdx := len(depth) - 1

			for i, d := range depth {
				if !d {
					b.WriteString(c.icons.indent)
					continue
				}

				if isFirst && i == lastIdx {
					b.WriteString(c.icons.lastIndent)
				} else {
					b.WriteString(c.icons.none)
				}
			}

			if isFirst {
				b.WriteString(c.taskPriorityColor(t).Sprint(c.taskStatusIcon(t)))
			} else {
				b.WriteString(c.icons.none)
			}

			return b.String()
		},
	}

	row := table.Row{t.ID, project, content, t.Description, strings.Join(labels, " "), due, deadline}

	rows := []table.Row{row}
	for i, st := range t.SubTasks {
		lastIdx := len(t.SubTasks) - 1
		rows = append(rows, c.renderTask(st, append(depth, i == lastIdx))...)
	}
	return rows
}

func (c *cli) PrintTasks(ts []*model.Task) {
	tbl := table.NewTable()
	tbl.AddHeader("ID", "Project", "Name", "Description", "Labels", "Due", "Deadline")
	tbl.SetHeaderStyle(&table.CellStyle{CellAttrs: text.Colors{text.Underline}})
	tbl.SetColStyle(2, &table.CellStyle{WrapText: table.BoolPtr(true)})
	tbl.SetColStyle(3, &table.CellStyle{WrapText: table.BoolPtr(true)})

	for _, t := range ts {
		tbl.AddRows(c.renderTask(t, []bool{}))
	}

	fmt.Print(tbl.Render())
}
