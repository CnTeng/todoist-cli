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
	"github.com/adrg/xdg"
	"github.com/urfave/cli/v3"
)

func newCmd() (*cli.Command, error) {
	appName := "todoist"
	configFilePath, err := xdg.ConfigFile(appName + "/config.json")
	if err != nil {
		return nil, fmt.Errorf("getting config file path: %v", err)
	}
	cfg := &model.Config{}

	return &cli.Command{
		Name:  "todoist",
		Usage: "A CLI for Todoist",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config",
				Usage:       "config file",
				Value:       configFilePath,
				Destination: &configFilePath,
			},
		},
		Before: func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
			file, err := os.ReadFile(configFilePath)
			if err != nil {
				return nil, fmt.Errorf("error reading config file: %v", err)
			}

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
	}, nil
}

func Execute() error {
	cmd, err := newCmd()
	if err != nil {
		return err
	}
	return cmd.Run(context.Background(), os.Args)
}
