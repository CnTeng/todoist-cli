package daemon

const (
	Sync string = "sync"

	// Task services
	TaskList     string = "Task.List"
	TaskAdd      string = "Task.Add"
	TaskQuickAdd string = "Task.QuickAdd"
	TaskModify   string = "Task.Modify"
	TaskMove     string = "Task.Move"
	TaskReorder  string = "Task.Reorder"
	TaskClose    string = "Task.Close"
	TaskReopen   string = "Task.Reopen"
	TaskRemove   string = "Task.Remove"

	// Project services
	ProjectList      string = "Project.List"
	ProjectAdd       string = "Project.Add"
	ProjectModify    string = "Project.Modify"
	ProjectReorder   string = "Project.Reorder"
	ProjectArchive   string = "Project.Archive"
	ProjectUnarchive string = "Project.Unarchive"
	ProjectRemove    string = "Project.Remove"

	// Section services
	SectionList      string = "Section.List"
	SectionAdd       string = "Section.Add"
	SectionModify    string = "Section.Modify"
	SectionMove      string = "Section.Move"
	SectionReorder   string = "Section.Reorder"
	SectionArchive   string = "Section.Archive"
	SectionUnarchive string = "Section.Unarchive"
	SectionRemove    string = "Section.Remove"

	// Label services
	LabelList    string = "Label.List"
	LabelAdd     string = "Label.Add"
	LabelModify  string = "Label.Modify"
	LabelReorder string = "Label.Reorder"
	LabelRemove  string = "Label.Remove"

	// Filter services
	FilterList    string = "Filter.List"
	FilterAdd     string = "Filter.Add"
	FilterModify  string = "Filter.Modify"
	FilterReorder string = "Filter.Reorder"
	FilterRemove  string = "Filter.Remove"
)
