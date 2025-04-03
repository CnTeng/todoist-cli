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

func NewRemoveCmd() *cli.Command {
	params := []string{}
	return &cli.Command{
		Name:        "remove",
		Aliases:     []string{"rm"},
		Usage:       "Remove a project",
		Description: "Remove a project in todoist",
		Arguments: []cli.Argument{
			&cli.StringArg{
				Name:   "id",
				Min:    1,
				Max:    -1,
				Values: &params,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			conn, err := net.Dial("unix", "@todo.sock")
			if err != nil {
				return err
			}
			defer conn.Close()

			cli := jrpc2.NewClient(channel.Line(conn, conn), nil)
			for _, id := range params {
				if _, err := cli.Call(ctx, daemon.ProjectRemove, &sync.ProjectDeleteArgs{ID: id}); err != nil {
					return err
				}
				fmt.Printf("Project deleted: %s\n", id)
			}

			return nil
		},
	}
}
