package daemon

const (
	Sync          string = "sync"
	CompletedGet  string = "Completed.Get"
	TaskGet       string = "Task.Get"
	TaskList      string = "Task.List"
	TaskAdd       string = "Task.Add"
	TaskQuickAdd  string = "Task.QuickAdd"
	TaskModify    string = "Task.Modify"
	TaskRemove    string = "Task.Remove"
	TaskClose     string = "Task.Close"
	TaskMove      string = "Task.Move"
	TaskReopen    string = "Task.Reopen"
	ProjectGet    string = "Project.Get"
	ProjectList   string = "Project.List"
	ProjectAdd    string = "Project.Add"
	ProjectModify string = "Project.Modify"
	ProjectRemove string = "Project.Remove"

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
