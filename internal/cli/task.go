package cli

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/CnTeng/table"
	"github.com/CnTeng/todoist-cli/internal/model"
	"github.com/CnTeng/todoist-cli/internal/utils"
	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/text"
)

type TaskListArgs struct {
	Completed   bool
	Tree        bool
	Description bool
}

func (c *cli) renderTask(t *model.Task, depth []bool, opt *TaskListArgs) []table.Row {
	row := table.Row{
		t.ID,
		c.taskProject(t),
		c.taskContent(t, depth, opt),
	}
	if opt.Description {
		row = append(row, t.Description)
	}
	row = append(row, c.taskLabels(t), c.taskDue(t), c.taskDeadline(t), c.taskDuration(t))

	rows := []table.Row{row}
	if !opt.Tree {
		return rows
	}

	for i, st := range t.SubTasks {
		lastIdx := len(t.SubTasks) - 1
		rows = append(rows, c.renderTask(st, append(depth, i == lastIdx), opt)...)
	}
	return rows
}

func (c *cli) PrintTasks(ts []*model.Task, opt *TaskListArgs) {
	tbl := table.NewTable()
	if opt.Description {
		tbl.AddHeader("ID", "Project", "Name", "Description", "Labels", "Due", "Deadline", "Duration")
	} else {
		tbl.AddHeader("ID", "Project", "Name", "Labels", "Due", "Deadline", "Duration")
	}
	tbl.SetHeaderStyle(&table.CellStyle{CellAttrs: text.Colors{text.Underline}})
	tbl.SetColStyle(2, &table.CellStyle{WrapText: utils.BoolPtr(true)})
	tbl.SetColStyle(3, &table.CellStyle{WrapText: utils.BoolPtr(true)})

	for _, t := range ts {
		tbl.AddRows(c.renderTask(t, []bool{}, opt))
	}

	fmt.Print(tbl.Render())
}

func (c *cli) taskProject(t *model.Task) string {
	return color.RGB(t.Project.Color.RGB()).Sprint(t.Project.Name)
}

func (c *cli) taskLabels(t *model.Task) string {
	labels := []string{}
	for _, l := range t.Labels {
		labels = append(labels, color.BgRGB(l.Color.RGB()).Sprint(l.Name))
	}
	return strings.Join(labels, " ")
}

func (c *cli) taskDue(t *model.Task) string {
	if t.Due != nil && t.Due.String != nil {
		return *t.Due.String
	}
	return ""
}

func (c *cli) taskDeadline(t *model.Task) string {
	if t.Deadline != nil {
		return t.Deadline.Date.Format(time.DateOnly)
	}
	return ""
}

func (c *cli) taskDuration(t *model.Task) string {
	if t.Duration != nil {
		return t.Duration.String()
	}
	return ""
}

func (c *cli) taskContent(t *model.Task, depth []bool, args *TaskListArgs) *table.Cell {
	depth = slices.Clone(depth)

	pColor := text.FgWhite
	switch t.Priority {
	case 1:
		pColor = text.FgWhite
	case 2:
		pColor = text.FgBlue
	case 3:
		pColor = text.FgYellow
	case 4:
		pColor = text.FgRed
	}

	sIcon := c.icons.undone
	if t.CompletedAt != nil {
		sIcon = c.icons.done
	}

	return &table.Cell{
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
				b.WriteString(pColor.Sprint(sIcon))
			} else {
				b.WriteString(c.icons.none)
			}

			if isFirst && !args.Tree && t.SubTaskStatus.Total > 0 {
				fmt.Fprintf(b, "(%d/%d) ", t.SubTaskStatus.Completed, t.SubTaskStatus.Total)
			}

			return b.String()
		},
	}
}
