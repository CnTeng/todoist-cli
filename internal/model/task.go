package model

import "github.com/CnTeng/todoist-api-go/sync"

type SubTaskStatus struct {
	Total     int `json:"total"`
	Completed int `json:"completed"`
}

type Task struct {
	*sync.Task
	Project       *sync.Project `json:"project"`
	SubTasks      []*Task       `json:"sub_tasks"`
	SubTaskStatus SubTaskStatus `json:"sub_task_status"`
	Labels        []*sync.Label `json:"labels"`
}
