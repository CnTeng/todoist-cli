package view

import (
	"github.com/CnTeng/table"
	"github.com/CnTeng/todoist-api-go/sync"
	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/text"
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
	if len(v.projects) == 0 {
		return "No projects found."
	}

	tbl := table.NewTable()
	tbl.AddHeader("ID", " ", "Name", "Color")
	tbl.SetHeaderStyle(headerStyle)
	tbl.SetColStyle(2, &table.CellStyle{
		WrapText: boolPtr(true),
		Markdown: boolPtr(true),
	})

	for _, p := range v.projects {
		row := table.Row{}

		row = append(row, p.ID)
		row = append(row, v.projectIcon(p))
		row = append(row, v.projectName(p))
		row = append(row, v.projectColor(p))

		tbl.AddRow(row)
	}

	return tbl.Render()
}

func (v *projectView) projectIcon(p *sync.Project) string {
	if p.InboxProject {
		return v.icons.Inbox
	} else if p.IsFavorite {
		return v.icons.Favorite
	}
	return v.icons.None
}

func (v *projectView) projectName(p *sync.Project) string {
	if p.IsArchived {
		return text.CrossedOut.Sprint(p.Name)
	}
	return p.Name
}

func (v *projectView) projectColor(p *sync.Project) string {
	return color.BgRGB(p.Color.RGB()).Sprint(p.Color)
}
