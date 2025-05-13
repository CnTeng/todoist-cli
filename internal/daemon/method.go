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

	// Label services
	LabelAdd     string = "Label.Add"
	LabelList    string = "Label.List"
	LabelModify  string = "Label.Modify"
	LabelRemove  string = "Label.Remove"
	LabelReorder string = "Label.Reorder"
)
