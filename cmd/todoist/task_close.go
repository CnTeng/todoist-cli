package cmd

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

var taskCloseArgs = []string{}

var taskCloseCmd = &cli.Command{
	Name:        "close",
	Aliases:     []string{"done"},
	Usage:       "Close a task",
	Description: "Close a task in todoist",
	Arguments: []cli.Argument{
		&cli.StringArg{
			Name:   "id",
			Min:    1,
			Max:    -1,
			Values: &taskCloseArgs,
		},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		conn, err := net.Dial("unix", "@todo.sock")
		if err != nil {
			return err
		}
		defer conn.Close()

		cli := jrpc2.NewClient(channel.Line(conn, conn), nil)
		for _, id := range taskCloseArgs {
			if _, err := cli.Call(ctx, daemon.TaskClose, &sync.ItemCloseArgs{ID: id}); err != nil {
				return err
			}
			fmt.Printf("Task close: %s\n", id)
		}

		return nil
	},
}
