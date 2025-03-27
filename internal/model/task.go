package model

import "github.com/CnTeng/todoist-api-go/sync/v9"

type Task struct {
	*sync.Item
	Project string        `json:"project"`
	Labels  []*sync.Label `json:"label"`
}
