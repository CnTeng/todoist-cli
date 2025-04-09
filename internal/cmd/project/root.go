package project

import (
	"github.com/CnTeng/todoist-cli/internal/cmd/util"
	"github.com/urfave/cli/v3"
)

func NewCmd(f *util.Factory) *cli.Command {
	return &cli.Command{
		Name:    "project",
		Aliases: []string{"proj"},
		Usage:   "project commands",
		Commands: []*cli.Command{
			NewListCmd(f),
			NewAddCmd(f),
			NewModifyCmd(f),
			NewRemoveCmd(f),
		},
	}
}
