package model

import "github.com/CnTeng/todoist-api-go/sync"

type Task struct {
	*sync.Item
	Project  *sync.Project `json:"project"`
	SubTasks []*Task       `json:"sub_tasks"`
	Labels   []*sync.Label `json:"labels"`
}
