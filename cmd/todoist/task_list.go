package cmd

import (
	"context"
	"net"

	tcli "github.com/CnTeng/todoist-cli/internal/cli"
	"github.com/CnTeng/todoist-cli/internal/daemon"
	"github.com/CnTeng/todoist-cli/internal/model"
	"github.com/creachadair/jrpc2"
	"github.com/creachadair/jrpc2/channel"
	"github.com/urfave/cli/v3"
)

var taskListArgs = &tcli.TaskListArgs{}

var taskListCmd = &cli.Command{
	Name:                   "list",
	Aliases:                []string{"ls"},
	Usage:                  "List tasks",
	Description:            "List tasks in todoist",
	UseShortOptionHandling: true,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:        "all",
			Aliases:     []string{"a"},
			Usage:       "list all tasks include completed",
			HideDefault: true,
			OnlyOnce:    true,
			Destination: &taskListArgs.Completed,
		},
		&cli.BoolFlag{
			Name:        "description",
			Aliases:     []string{"d"},
			Usage:       "list tasks include description",
			HideDefault: true,
			OnlyOnce:    true,
			Destination: &taskListArgs.Description,
		},
		&cli.BoolFlag{
			Name:        "tree",
			Aliases:     []string{"t"},
			Usage:       "list tasks in tree format",
			HideDefault: true,
			OnlyOnce:    true,
			Destination: &taskListArgs.Tree,
		},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		conn, err := net.Dial("unix", "@todo.sock")
		if err != nil {
			return err
		}
		defer conn.Close()

		cli := jrpc2.NewClient(channel.Line(conn, conn), nil)
		result := []*model.Task{}
		if err := cli.CallResult(
			ctx,
			daemon.TaskList,
			struct {
				All bool `json:"all"`
			}{All: taskListArgs.Completed},
			&result,
		); err != nil {
			return err
		}

		c := tcli.NewCLI(tcli.Nerd)
		c.PrintTasks(result, taskListArgs)

		return nil
	},
}
