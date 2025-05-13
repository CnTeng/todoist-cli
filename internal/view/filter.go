package view

import (
	"github.com/CnTeng/table"
	"github.com/CnTeng/todoist-api-go/sync"
	"github.com/fatih/color"
)

type filterView struct {
	icons   *Icons
	filters []*sync.Filter
}

func NewFilterView(filters []*sync.Filter, icons *Icons) View {
	return &filterView{
		icons:   icons,
		filters: filters,
	}
}

func (v *filterView) Render() string {
	if len(v.filters) == 0 {
		return "No filters found."
	}

	tbl := table.NewTable()
	tbl.AddHeader("ID", " ", "Name", "Color")
	tbl.SetHeaderStyle(headerStyle)

	for _, f := range v.filters {
		row := table.Row{
			f.ID,
			v.filterIcon(f),
			f.Name,
			v.labelColor(f),
		}
		tbl.AddRow(row)
	}

	return tbl.Render()
}

func (v *filterView) filterIcon(f *sync.Filter) string {
	if f.IsFavorite {
		return v.icons.Favorite
	}
	return v.icons.None
}

func (v *filterView) labelColor(f *sync.Filter) string {
	return color.BgRGB(f.Color.RGB()).Sprint(f.Color)
}
