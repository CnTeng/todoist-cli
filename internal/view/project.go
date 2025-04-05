package view

import (
	"github.com/CnTeng/table"
	"github.com/CnTeng/todoist-api-go/sync"
	"github.com/fatih/color"
)

type projectView struct {
	icons    *Icons
	projects []*sync.Project
}

func NewProjectView(projects []*sync.Project, icons *Icons) View {
	return &projectView{
		icons:    icons,
		projects: projects,
	}
}

func (v *projectView) Render() string {
	tbl := table.NewTable()
	tbl.AddHeader("ID", "Name", "Color")
	tbl.SetHeaderStyle(headerStyle)
	tbl.SetColStyle(1, &table.CellStyle{WrapText: boolPtr(true)})

	for _, p := range v.projects {
		row := table.Row{}

		row = append(row, p.ID)
		row = append(row, v.projectName(p))
		row = append(row, color.BgRGB(p.Color.RGB()).Sprint(p.Color))

		tbl.AddRow(row)
	}

	return tbl.Render()
}

func (v *projectView) projectName(p *sync.Project) *table.Cell {
	icon := v.icons.none
	if p.InboxProject {
		icon = v.icons.inbox
	} else if p.IsFavorite {
		icon = v.icons.favorite
	}
	return &table.Cell{
		Content: p.Name,
		PrefixFunc: func(isFirst, isLast bool) string {
			if isFirst {
				return icon
			}
			return v.icons.none
		},
	}
}
