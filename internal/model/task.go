package model

import "github.com/CnTeng/todoist-api-go/sync"

type SubTaskStatus struct {
	Total     int `json:"total"`
	Completed int `json:"completed"`
}

type Task struct {
	*sync.Task
	ProjectName   string        `json:"project_name"`
	ProjectColor  sync.Color    `json:"project_color"`
	SectionName   string        `json:"section_name"`
	SubTasks      []*Task       `json:"sub_tasks"`
	SubTaskStatus SubTaskStatus `json:"sub_task_status"`
	Labels        []*Label      `json:"labels"`
}

type TaskListArgs struct {
	Completed bool   `json:"completed"`
	ProjectID string `json:"project_id"`
}
