package view

import (
	"github.com/CnTeng/table"
	"github.com/CnTeng/todoist-cli/internal/model"
	"github.com/fatih/color"
)

type labelView struct {
	icons  *Icons
	labels []*model.Label
}

func NewLabelView(labels []*model.Label, icons *Icons) View {
	return &labelView{
		icons:  icons,
		labels: labels,
	}
}

func (v *labelView) Render() string {
	if len(v.labels) == 0 {
		return "No labels found."
	}

	tbl := table.NewTable()
	tbl.AddHeader("ID", " ", "Name", "Color", "Tasks Count")
	tbl.SetHeaderStyle(headerStyle)

	for _, l := range v.labels {
		row := table.Row{
			v.labelID(l),
			v.labelIcon(l),
			l.Name,
			v.labelColor(l),
			l.Count,
		}
		tbl.AddRow(row)
	}

	return tbl.Render()
}

func (v *labelView) labelID(l *model.Label) string {
	if l.IsShared {
		return ""
	}
	return l.ID
}

func (v *labelView) labelIcon(l *model.Label) string {
	if l.IsFavorite {
		return v.icons.Favorite
	}
	return v.icons.None
}

func (v *labelView) labelColor(l *model.Label) string {
	if l.IsShared {
		return "shared"
	}
	return color.BgRGB(l.Color.RGB()).Sprint(l.Color)
}
