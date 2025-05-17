package model

import "github.com/CnTeng/todoist-api-go/sync"

type Section struct {
	*sync.Section
	ProjectName string `json:"project_name"`
}
