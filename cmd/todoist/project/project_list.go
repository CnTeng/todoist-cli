package project

import (
	"context"

	"github.com/CnTeng/todoist-api-go/sync"
	"github.com/CnTeng/todoist-cli/cmd/todoist/util"
	tcli "github.com/CnTeng/todoist-cli/internal/cli"
	"github.com/CnTeng/todoist-cli/internal/daemon"
	"github.com/urfave/cli/v3"
)

func NewListCmd(f *util.Factory) *cli.Command {
	return &cli.Command{
		Name:    "list",
		Aliases: []string{"ls"},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			result := []*sync.Project{}
			if err := f.RpcClient.CallResult(ctx, daemon.ProjectList, nil, &result); err != nil {
				return err
			}

			c := tcli.NewCLI(tcli.Nerd)
			c.PrintProjects(result)

			return nil
		},
	}
}
