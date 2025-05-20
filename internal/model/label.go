package model

import "github.com/CnTeng/todoist-api-go/sync"

type Label struct {
	*sync.Label
	IsShared bool `json:"is_shared"`
	Count    int  `json:"count"`
}

type LabelUpdateArgs struct {
	Name string `json:"name"`
	Args sync.LabelUpdateArgs
}

type LabelDeleteArgs struct {
	Name string `json:"name"`
}
