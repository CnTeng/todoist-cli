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

var taskMoveArgs = &sync.ItemMoveArgs{}

var taskMoveCmd = &cli.Command{
	Name:        "move",
	Aliases:     []string{"mv"},
	Usage:       "Move a task",
	Description: "Move a task in todoist",
	Arguments: []cli.Argument{
		&cli.StringArg{
			Name:        "id",
			Min:         1,
			Max:         1,
			Destination: &taskMoveArgs.ID,
		},
	},
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "section",
			Aliases:  []string{"s"},
			Usage:    "Section ID",
			OnlyOnce: true,
			Action: func(_ context.Context, _ *cli.Command, v string) error {
				taskMoveArgs.SectionID = &v
				return nil
			},
		},
		&cli.StringFlag{
			Name:     "parent",
			Usage:    "Parent task ID",
			OnlyOnce: true,
			Action: func(_ context.Context, _ *cli.Command, v string) error {
				taskMoveArgs.ParentID = &v
				return nil
			},
		},
		&cli.StringFlag{
			Name:     "project",
			Aliases:  []string{"P"},
			Usage:    "Project ID",
			OnlyOnce: true,
			Action: func(_ context.Context, _ *cli.Command, v string) error {
				taskMoveArgs.ProjectID = &v
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
		if _, err := cli.Call(ctx, daemon.TaskMove, taskMoveArgs); err != nil {
			return err
		}

		fmt.Printf("Task moved: %s\n", taskMoveArgs.ID)

		return nil
	},
}
