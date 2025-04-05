package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/CnTeng/todoist-cli/cmd/todoist/daemon"
	"github.com/CnTeng/todoist-cli/cmd/todoist/project"
	"github.com/CnTeng/todoist-cli/cmd/todoist/sync"
	"github.com/CnTeng/todoist-cli/cmd/todoist/task"
	"github.com/CnTeng/todoist-cli/cmd/todoist/util"
	"github.com/CnTeng/todoist-cli/internal/model"
	"github.com/adrg/xdg"
	"github.com/creachadair/jrpc2"
	"github.com/creachadair/jrpc2/channel"
	"github.com/urfave/cli/v3"
)

func newCmd() (*cli.Command, error) {
	appName := "todoist"
	configFilePath, err := xdg.ConfigFile(appName + "/config.json")
	if err != nil {
		return nil, fmt.Errorf("getting config file path: %v", err)
	}

	f := &util.Factory{
		Config: &model.Config{},
	}

	var conn net.Conn

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

			if err := json.Unmarshal(file, f.Config); err != nil {
				return nil, fmt.Errorf("error decoding config file: %v", err)
			}

			if cmd.Args().First() == "daemon" {
				return ctx, nil
			}
			conn, err = net.Dial("unix", "@todo.sock")
			if err != nil {
				return nil, err
			}

			f.RpcClient = jrpc2.NewClient(channel.Line(conn, conn), nil)

			return ctx, nil
		},

		After: func(ctx context.Context, cmd *cli.Command) error {
			if conn == nil {
				return nil
			}
			if err := conn.Close(); err != nil {
				return fmt.Errorf("closing connection: %v", err)
			}
			return nil
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
