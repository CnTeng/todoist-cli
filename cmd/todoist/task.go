package cmd

import (
	"context"
	"fmt"
	"net"

	"github.com/CnTeng/todoist-api-go/sync"
	tcli "github.com/CnTeng/todoist-cli/internal/cli"
	"github.com/CnTeng/todoist-cli/internal/daemon"
	"github.com/CnTeng/todoist-cli/internal/model"
	"github.com/creachadair/jrpc2"
	"github.com/creachadair/jrpc2/channel"
	"github.com/urfave/cli/v3"
)

var taskListCmd = &cli.Command{
	Name:    "list",
	Aliases: []string{"ls"},
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:     "tree",
			Aliases:  []string{"t"},
			Usage:    "list tasks in tree format",
			OnlyOnce: true,
			Value:    false,
		},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		conn, err := net.Dial("unix", "@todo.sock")
		if err != nil {
			fmt.Printf("Error dialing daemon: %v\n", err)
		}
		defer conn.Close()
		cli := jrpc2.NewClient(channel.Line(conn, conn), nil)

		resp := []*model.Task{}
		if err := cli.CallResult(ctx, "listTasks", nil, &resp); err != nil {
			fmt.Printf("Error calling taskLists: %v\n", err)
		}

		c := tcli.NewCLI(tcli.Nerd)
		c.PrintTasks(resp, cmd.Bool("tree"))

		return nil
	},
}

var taskAddCmd = &cli.Command{
	Name:    "add",
	Aliases: []string{"a"},
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "content",
			Aliases:  []string{"c"},
			Usage:    "Task content",
			Required: true,
			OnlyOnce: true,
		},
		&cli.StringFlag{
			Name:     "description",
			Aliases:  []string{"d"},
			Usage:    "Task description",
			OnlyOnce: true,
		},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		conn, err := net.Dial("unix", "@todo.sock")
		if err != nil {
			fmt.Printf("Error dialing daemon: %v\n", err)
		}
		defer conn.Close()
		cli := jrpc2.NewClient(channel.Line(conn, conn), nil)

		args := &sync.ItemAddArgs{}
		args.Content = cmd.String("content")
		if cmd.IsSet("description") {
			description := cmd.String("description")
			args.Description = &description
		}

		if _, err := cli.Call(ctx, daemon.AddTask, args); err != nil {
			fmt.Printf("Error calling add task: %v\n", err)
		}

		resp := []*model.Task{}
		if err := cli.CallResult(ctx, daemon.GetTask, nil, &resp); err != nil {
			fmt.Printf("Error calling taskLists: %v\n", err)
		}

		c := tcli.NewCLI(tcli.Nerd)
		c.PrintTasks(resp, false)

		return nil
	},
}
var taskRemoveArgs = []string{}

var taskDeleteCmd = &cli.Command{
	Name:    "remove",
	Aliases: []string{"rm"},
	Arguments: []cli.Argument{
		&cli.StringArg{
			Name:   "id",
			Values: &taskRemoveArgs,
			Min:    1,
			Max:    -1,
		},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		conn, err := net.Dial("unix", "@todo.sock")
		if err != nil {
			fmt.Printf("Error dialing daemon: %v\n", err)
		}
		defer conn.Close()
		cli := jrpc2.NewClient(channel.Line(conn, conn), nil)

		for _, id := range taskRemoveArgs {
			if _, err := cli.Call(ctx, daemon.TaskDelete, &sync.ItemDeleteArgs{ID: id}); err != nil {
				fmt.Printf("Error calling delete task: %v\n", err)
			}
		}

		// resp := []*model.Task{}
		// if err := cli.CallResult(ctx, daemon.GetTask, nil, &resp); err != nil {
		// 	fmt.Printf("Error calling taskLists: %v\n", err)
		// }
		//
		// c := tcli.NewCLI(tcli.Nerd)
		// c.PrintTasks(resp, false)

		return nil
	},
}
