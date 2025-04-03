package project

import (
	"context"
	"fmt"
	"net"

	"github.com/CnTeng/todoist-api-go/sync"
	"github.com/CnTeng/todoist-cli/internal/daemon"
	"github.com/creachadair/jrpc2"
	"github.com/creachadair/jrpc2/channel"
	"github.com/urfave/cli/v3"
)

func NewAddCmd() *cli.Command {
	params := &sync.ProjectAddArgs{}
	return &cli.Command{
		Name:        "add",
		Aliases:     []string{"a"},
		Usage:       "Add a project",
		Description: "Add a project to todoist",
		Arguments: []cli.Argument{
			&cli.StringArg{
				Name:        "name",
				Min:         1,
				Max:         1,
				Destination: &params.Name,
				Config:      cli.StringConfig{TrimSpace: true},
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "color",
				Aliases:  []string{"c"},
				Usage:    "Task color",
				OnlyOnce: true,
				Action: func(_ context.Context, _ *cli.Command, v string) error {
					color, err := sync.ParseColor(v)
					if err != nil {
						return err
					}
					params.Color = &color
					return nil
				},
			},
			&cli.BoolFlag{
				Name:    "favorite",
				Aliases: []string{"f"},
				Usage:   "Add project to favorites",
				Action: func(_ context.Context, _ *cli.Command, v bool) error {
					params.IsFavorite = &v
					return nil
				},
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			conn, err := net.Dial("unix", "@todo.sock")
			if err != nil {
				return err
			}
			defer conn.Close()

			cli := jrpc2.NewClient(channel.Line(conn, conn), nil)
			if _, err := cli.Call(ctx, daemon.ProjectAdd, params); err != nil {
				return err
			}

			fmt.Printf("Project added: %s\n", params.Name)

			return nil
		},
	}
}
