package view

import (
	"github.com/CnTeng/table"
	"github.com/CnTeng/todoist-cli/internal/model"
	"github.com/jedib0t/go-pretty/v6/text"
)

type sectionView struct {
	icons    *Icons
	sections []*model.Section
}

func NewSectionView(sections []*model.Section, icons *Icons) *sectionView {
	return &sectionView{
		icons:    icons,
		sections: sections,
	}
}

func (v *sectionView) Render() string {
	if len(v.sections) == 0 {
		return "No sections found."
	}

	tbl := table.NewTable()
	tbl.AddHeader("ID", "Project", "Name")
	tbl.SetHeaderStyle(headerStyle)

	for _, s := range v.sections {
		row := table.Row{
			s.ID,
			s.ProjectName,
			v.sectionName(s),
		}
		tbl.AddRow(row)
	}

	return tbl.Render()
}

func (v *sectionView) sectionName(s *model.Section) string {
	if s.IsArchived {
		return text.CrossedOut.Sprint(s.Name)
	}
	return s.Name
}
