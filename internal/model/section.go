package model

import "github.com/CnTeng/todoist-api-go/sync"

type Section struct {
	*sync.Section
	ProjectName string `json:"project_name"`
}

type SectionListArgs struct {
	ProjectID    string `json:"project_id"`
	All          bool   `json:"all"`
	OnlyArchived bool   `json:"only_archived"`
}
