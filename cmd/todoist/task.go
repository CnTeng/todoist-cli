package cmd

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/CnTeng/todoist-api-go/sync"
	tcli "github.com/CnTeng/todoist-cli/internal/cli"
	"github.com/CnTeng/todoist-cli/internal/daemon"
	"github.com/CnTeng/todoist-cli/internal/model"
	"github.com/CnTeng/todoist-cli/internal/utils"
	"github.com/creachadair/jrpc2"
	"github.com/creachadair/jrpc2/channel"
	"github.com/urfave/cli/v3"
)

var (
	taskListArgs   = &tcli.TaskListArgs{}
	taskAddArgs    = &sync.ItemAddArgs{}
	taskRemoveArgs = []string{}
)

var taskListCmd = &cli.Command{
	Name:    "list",
	Aliases: []string{"ls"},
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:        "description",
			Aliases:     []string{"d"},
			Usage:       "list tasks include description",
			OnlyOnce:    true,
			Value:       false,
			Destination: &taskListArgs.Description,
		},
		&cli.BoolFlag{
			Name:        "tree",
			Aliases:     []string{"t"},
			Usage:       "list tasks in tree format",
			OnlyOnce:    true,
			Value:       false,
			Destination: &taskListArgs.Tree,
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
		c.PrintTasks(resp, taskListArgs)

		return nil
	},
}

var taskAddCmd = &cli.Command{
	Name:                  "add",
	Aliases:               []string{"a"},
	EnableShellCompletion: true,
	Arguments: []cli.Argument{
		&cli.StringArg{
			Name:        "content",
			Destination: &taskAddArgs.Content,
			Min:         1,
			Max:         1,
			Config:      cli.StringConfig{TrimSpace: true},
		},
	},
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "description",
			Aliases:  []string{"D"},
			Usage:    "Task description",
			OnlyOnce: true,
		},
		&cli.StringFlag{
			Name:     "project",
			Aliases:  []string{"P"},
			Usage:    "Project ID",
			OnlyOnce: true,
		},
		&cli.StringFlag{
			Name:     "due",
			Aliases:  []string{"d"},
			Usage:    "Due date",
			OnlyOnce: true,
		},
		&cli.TimestampFlag{
			Name:     "deadline",
			Usage:    "Deadline date",
			OnlyOnce: true,
			Config: cli.TimestampConfig{
				Layouts: []string{time.DateOnly},
			},
		},
		&cli.IntFlag{
			Name:     "priority",
			Aliases:  []string{"p"},
			Usage:    "Task priority",
			OnlyOnce: true,
		},
		&cli.StringFlag{
			Name:     "parent",
			Usage:    "Parent task ID",
			OnlyOnce: true,
		},
		&cli.StringFlag{
			Name:     "section",
			Aliases:  []string{"s"},
			Usage:    "Section ID",
			OnlyOnce: true,
		},
		&cli.StringSliceFlag{
			Name:     "labels",
			Aliases:  []string{"l"},
			Usage:    "Labels",
			OnlyOnce: true,
		},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		conn, err := net.Dial("unix", "@todo.sock")
		if err != nil {
			return err
		}
		defer func() { _ = conn.Close() }()

		cli := jrpc2.NewClient(channel.Line(conn, conn), nil)

		if cmd.IsSet("description") {
			taskAddArgs.Description = utils.StringPtr(cmd.String("description"))
		}
		if cmd.IsSet("project") {
			taskAddArgs.ProjectID = utils.StringPtr(cmd.String("project"))
		}
		if cmd.IsSet("due") {
			taskAddArgs.Due = &sync.Due{String: utils.StringPtr(cmd.String("due"))}
		}
		if cmd.IsSet("deadline") {
			taskAddArgs.Deadline = &sync.Deadline{Date: cmd.Timestamp("deadline"), Lang: cfg.Lang}
		}
		if cmd.IsSet("priority") {
			taskAddArgs.Priority = utils.IntPtr(int(cmd.Int("priority")))
		}
		if cmd.IsSet("parent") {
			taskAddArgs.ParentID = utils.StringPtr(cmd.String("parent"))
		}
		if cmd.IsSet("section") {
			taskAddArgs.SectionID = utils.StringPtr(cmd.String("section"))
		}
		if cmd.IsSet("labels") {
			taskAddArgs.Labels = cmd.StringSlice("labels")
		}

		if _, err := cli.Call(ctx, daemon.AddTask, taskAddArgs); err != nil {
			fmt.Printf("Error calling add task: %v\n", err)
		}

		fmt.Printf("Task added: %s\n", taskAddArgs.Content)

		return nil
	},
}

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
