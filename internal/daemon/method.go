package daemon

const (
	Sync         string = "sync"
	CompletedGet string = "Completed.Get"

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
	ProjectAdd       string = "Project.Add"
	ProjectList      string = "Project.List"
	ProjectModify    string = "Project.Modify"
	ProjectReorder   string = "Project.Reorder"
	ProjectArchive   string = "Project.Archive"
	ProjectUnarchive string = "Project.Unarchive"
	ProjectRemove    string = "Project.Remove"

	// Section services
	SectionAdd       string = "Section.Add"
	SectionList      string = "Section.List"
	SectionModify    string = "Section.Modify"
	SectionMove      string = "Section.Move"
	SectionArchive   string = "Section.Archive"
	SectionUnarchive string = "Section.Unarchive"
	SectionRemove    string = "Section.Remove"
	SectionReorder   string = "Section.Reorder"

	// Label services
	LabelAdd     string = "Label.Add"
	LabelList    string = "Label.List"
	LabelModify  string = "Label.Modify"
	LabelRemove  string = "Label.Remove"
	LabelReorder string = "Label.Reorder"

	// Filter services
	FilterAdd     string = "Filter.Add"
	FilterList    string = "Filter.List"
	FilterModify  string = "Filter.Modify"
	FilterRemove  string = "Filter.Remove"
	FilterReorder string = "Filter.Reorder"
)
