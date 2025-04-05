package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/CnTeng/todoist-cli/internal/cmd/daemon"
	"github.com/CnTeng/todoist-cli/internal/cmd/project"
	"github.com/CnTeng/todoist-cli/internal/cmd/sync"
	"github.com/CnTeng/todoist-cli/internal/cmd/task"
	"github.com/CnTeng/todoist-cli/internal/cmd/util"
	"github.com/adrg/xdg"
	"github.com/urfave/cli/v3"
)

func newCmd() (*cli.Command, error) {
	appName := "todoist"
	configFilePath, err := xdg.ConfigFile(appName + "/config.toml")
	if err != nil {
		return nil, fmt.Errorf("getting config file path: %v", err)
	}

	f := util.NewFactory()
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

			if err := toml.Unmarshal(file, f); err != nil {
				return nil, fmt.Errorf("error decoding config file: %v", err)
			}

			if cmd.Args().First() == "daemon" {
				return ctx, nil
			}

			if err := f.Dial(); err != nil {
				return nil, err
			}

			return ctx, nil
		},

		After: func(ctx context.Context, cmd *cli.Command) error {
			return f.Close()
		},

		Commands: []*cli.Command{
			task.NewListCmd(f),
			task.NewAddCmd(f),
			task.NewQuickAddCmd(f),
			task.NewModifyCmd(f),
			task.NewCloseCmd(f),
			task.NewReopenCmd(f),
			task.NewRemoveCmd(f),
			task.NewMoveCmd(f),
			project.NewCmd(f),
			sync.NewCmd(f),
			daemon.NewCmd(f),
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
