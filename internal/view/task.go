package view

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/CnTeng/table"
	"github.com/CnTeng/todoist-cli/internal/model"
	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/text"
)

type TaskViewConfig struct {
	Completed   bool `json:"completed"`
	Tree        bool `json:"-"`
	Description bool `json:"-"`
}

type taskView struct {
	icons  *Icons
	tasks  []*model.Task
	config *TaskViewConfig
}

func NewTaskView(tasks []*model.Task, icons *Icons, config *TaskViewConfig) View {
	return &taskView{
		icons:  icons,
		tasks:  tasks,
		config: config,
	}
}

func (v *taskView) Render() string {
	tbl := table.NewTable()
	if v.config.Description {
		tbl.AddHeader("ID", "Project", "Name", "Description", "Labels", "Due", "Deadline", "Duration")
	} else {
		tbl.AddHeader("ID", "Project", "Name", "Labels", "Due", "Deadline", "Duration")
	}
	tbl.SetHeaderStyle(headerStyle)
	tbl.SetColStyle(2, &table.CellStyle{
		WrapText: boolPtr(true),
		Markdown: boolPtr(true),
	})
	tbl.SetColStyle(3, &table.CellStyle{
		WrapText: boolPtr(true),
		Markdown: boolPtr(true),
	})

	for _, t := range v.tasks {
		tbl.AddRows(v.renderTask(t, []bool{}))
	}

	return tbl.Render()
}

func (v *taskView) renderTask(t *model.Task, depth []bool) []table.Row {
	row := table.Row{
		t.ID,
		v.taskProject(t),
		v.taskContent(t, depth),
	}
	if v.config.Description {
		row = append(row, t.Description)
	}
	row = append(row, v.taskLabels(t), v.taskDue(t), v.taskDeadline(t), v.taskDuration(t))

	rows := []table.Row{row}
	if !v.config.Tree {
		return rows
	}

	for i, st := range t.SubTasks {
		lastIdx := len(t.SubTasks) - 1
		rows = append(rows, v.renderTask(st, append(depth, i == lastIdx))...)
	}
	return rows
}

func (v *taskView) taskProject(t *model.Task) string {
	return color.RGB(t.Project.Color.RGB()).Sprint(t.Project.Name)
}

func (v *taskView) taskLabels(t *model.Task) string {
	labels := []string{}
	for _, l := range t.Labels {
		labels = append(labels, color.BgRGB(l.Color.RGB()).Sprint(l.Name))
	}
	return strings.Join(labels, " ")
}

func (v *taskView) taskDue(t *model.Task) string {
	if t.Due != nil && t.Due.String != nil {
		return *t.Due.String
	}
	return ""
}

func (v *taskView) taskDeadline(t *model.Task) string {
	if t.Deadline != nil {
		return t.Deadline.Date.Format(time.DateOnly)
	}
	return ""
}

func (v *taskView) taskDuration(t *model.Task) string {
	if t.Duration != nil {
		return t.Duration.String()
	}
	return ""
}

func (v *taskView) taskContent(t *model.Task, depth []bool) *table.Cell {
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

	sIcon := v.icons.Undone
	if t.Checked {
		sIcon = v.icons.Done
	}

	return &table.Cell{
		Content: t.Content,
		PrefixFunc: func(isFirst, isLast bool) string {
			b := &strings.Builder{}
			lastIdx := len(depth) - 1

			for i, d := range depth {
				if !d {
					b.WriteString(v.icons.Indent)
					continue
				}

				if isFirst && i == lastIdx {
					b.WriteString(v.icons.LastIndent)
				} else {
					b.WriteString(v.icons.None)
				}
			}

			if isFirst {
				b.WriteString(pColor.Sprint(sIcon))
			} else {
				b.WriteString(v.icons.None)
			}

			if isFirst && !v.config.Tree && t.SubTaskStatus.Total > 0 {
				fmt.Fprintf(b, "(%d/%d) ", t.SubTaskStatus.Completed, t.SubTaskStatus.Total)
			}

			return b.String()
		},
	}
}
