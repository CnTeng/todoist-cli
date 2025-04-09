package cmd

import (
	"context"
	"os"
	"path/filepath"

	"github.com/CnTeng/todoist-cli/internal/cmd/daemon"
	"github.com/CnTeng/todoist-cli/internal/cmd/project"
	"github.com/CnTeng/todoist-cli/internal/cmd/sync"
	"github.com/CnTeng/todoist-cli/internal/cmd/task"
	"github.com/CnTeng/todoist-cli/internal/cmd/util"
	"github.com/adrg/xdg"
	"github.com/urfave/cli/v3"
)

const (
	appName    = "todoist"
	configFile = "config.toml"
	cacheFile  = "todoist.db"
)

func newCmd() (*cli.Command, error) {
	configFilePath, err := xdg.ConfigFile(filepath.Join(appName, configFile))
	if err != nil {
		return nil, err
	}
	dataFilePath, err := xdg.DataFile(filepath.Join(appName, cacheFile))
	if err != nil {
		return nil, err
	}

	f := util.NewFactory(configFilePath, dataFilePath)

	return &cli.Command{
		Name:  appName,
		Usage: "A CLI for Todoist",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config",
				Usage:       "config file",
				Value:       f.ConfigFilePath,
				Destination: &f.ConfigFilePath,
			},
		},
		Before: func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
			if err := f.ReadConfig(); err != nil {
				return nil, err
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
