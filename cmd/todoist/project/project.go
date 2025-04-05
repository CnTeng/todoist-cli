package project

import (
	"github.com/CnTeng/todoist-cli/cmd/todoist/util"
	"github.com/urfave/cli/v3"
)

func NewCmd(f *util.Factory) *cli.Command {
	return &cli.Command{
		Name:  "project",
		Usage: "project commands",
		Commands: []*cli.Command{
			NewListCmd(f),
			NewAddCmd(f),
			NewRemoveCmd(f),
		},
	}
}
