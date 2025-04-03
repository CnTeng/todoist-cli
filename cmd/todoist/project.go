package cmd

import "github.com/urfave/cli/v3"

func projectCmd() *cli.Command {
	return &cli.Command{
		Name:  "project",
		Usage: "project commands",
		Commands: []*cli.Command{
			projectListCmd(),
		},
	}
}
