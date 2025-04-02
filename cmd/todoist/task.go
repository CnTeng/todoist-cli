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
	taskModifyArgs = &sync.ItemUpdateArgs{}
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
			return err
		}
		defer func() { _ = conn.Close() }()

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

var taskModifyCmd = &cli.Command{
	Name:                  "modify",
	Aliases:               []string{"m"},
	EnableShellCompletion: true,
	Arguments: []cli.Argument{
		&cli.StringArg{
			Name:        "ID",
			Destination: &taskModifyArgs.ID,
			Min:         1,
			Max:         1,
			Config:      cli.StringConfig{TrimSpace: true},
		},
	},
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "content",
			Aliases:  []string{"c"},
			Usage:    "Task content",
			OnlyOnce: true,
		},
		&cli.StringFlag{
			Name:     "description",
			Aliases:  []string{"D"},
			Usage:    "Task description",
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
		&cli.StringSliceFlag{
			Name:     "labels",
			Aliases:  []string{"l"},
			Usage:    "Labels",
			OnlyOnce: true,
		},
		&cli.StringFlag{
			Name:     "duration",
			Usage:    "Duration",
			OnlyOnce: true,
			Validator: func(v string) error {
				_, err := sync.ParseDuration(v)
				return err
			},
		},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		conn, err := net.Dial("unix", "@todo.sock")
		if err != nil {
			return err
		}
		defer func() { _ = conn.Close() }()

		cli := jrpc2.NewClient(channel.Line(conn, conn), nil)

		if cmd.IsSet("content") {
			taskModifyArgs.Content = utils.StringPtr(cmd.String("content"))
		}
		if cmd.IsSet("description") {
			taskModifyArgs.Description = utils.StringPtr(cmd.String("description"))
		}
		if cmd.IsSet("due") {
			taskModifyArgs.Due = &sync.Due{String: utils.StringPtr(cmd.String("due"))}
		}
		if cmd.IsSet("deadline") {
			taskModifyArgs.Deadline = &sync.Deadline{Date: cmd.Timestamp("deadline"), Lang: cfg.Lang}
		}
		if cmd.IsSet("priority") {
			taskModifyArgs.Priority = utils.IntPtr(int(cmd.Int("priority")))
		}
		if cmd.IsSet("labels") {
			taskModifyArgs.Labels = cmd.StringSlice("labels")
		}
		if cmd.IsSet("duration") {
			duration, err := sync.ParseDuration(cmd.String("duration"))
			if err != nil {
				fmt.Printf("Error parsing duration: %v\n", err)
			}
			taskModifyArgs.Duration = duration
		}

		if _, err := cli.Call(ctx, daemon.TaskModify, taskModifyArgs); err != nil {
			fmt.Printf("Error calling add task: %v\n", err)
		}

		fmt.Printf("Task modified: %s\n", taskModifyArgs.ID)

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
			return err
		}
		defer func() { _ = conn.Close() }()

		cli := jrpc2.NewClient(channel.Line(conn, conn), nil)

		for _, id := range taskRemoveArgs {
			if _, err := cli.Call(ctx, daemon.TaskDelete, &sync.ItemDeleteArgs{ID: id}); err != nil {
				fmt.Printf("Error calling delete task: %v\n", err)
			}
			fmt.Printf("Task deleted: %s\n", id)
		}

		return nil
	},
}
