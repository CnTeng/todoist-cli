package model

import "github.com/CnTeng/todoist-api-go/sync"

type Label struct {
	*sync.Label
	IsShared bool `json:"is_shared"`
	Count    int  `json:"count"`
}
