package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/CnTeng/todoist-cli/cmd/todoist/daemon"
	"github.com/CnTeng/todoist-cli/cmd/todoist/project"
	"github.com/CnTeng/todoist-cli/cmd/todoist/sync"
	"github.com/CnTeng/todoist-cli/cmd/todoist/task"
	"github.com/CnTeng/todoist-cli/internal/model"
	"github.com/urfave/cli/v3"
)

func newCmd() *cli.Command {
	var configFile string
	cfg := &model.Config{}

	return &cli.Command{
		Name:  "todoist",
		Usage: "A CLI for Todoist",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config",
				Usage:       "config file",
				Value:       "config.json",
				Destination: &configFile,
			},
		},
		Before: func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
			file, err := os.ReadFile(configFile)
			if err != nil {
				return nil, fmt.Errorf("error reading config file: %v", err)
			}

			cfg = &model.Config{}
			if err := json.Unmarshal(file, cfg); err != nil {
				return nil, fmt.Errorf("error decoding config file: %v", err)
			}

			return ctx, nil
		},

		Commands: []*cli.Command{
			task.NewListCmd(),
			task.NewAddCmd(cfg),
			task.NewQuickAddCmd(),
			task.NewModifyCmd(cfg),
			task.NewCloseCmd(),
			task.NewReopenCmd(),
			task.NewRemoveCmd(),
			task.NewMoveCmd(),
			project.NewCmd(),
			sync.NewCmd(),
			daemon.NewCmd(cfg),
		},
	}
}

func Execute() error {
	return newCmd().Run(context.Background(), os.Args)
}
