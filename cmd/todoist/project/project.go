package project

import "github.com/urfave/cli/v3"

func NewCmd() *cli.Command {
	return &cli.Command{
		Name:  "project",
		Usage: "project commands",
		Commands: []*cli.Command{
			NewListCmd(),
		},
	}
}
