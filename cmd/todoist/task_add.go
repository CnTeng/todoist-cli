package cmd

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/CnTeng/todoist-api-go/sync"
	"github.com/CnTeng/todoist-cli/internal/daemon"
	"github.com/CnTeng/todoist-cli/internal/utils"
	"github.com/creachadair/jrpc2"
	"github.com/creachadair/jrpc2/channel"
	"github.com/urfave/cli/v3"
)

var taskAddArgs = &sync.ItemAddArgs{}

var taskAddCmd = &cli.Command{
	Name:        "add",
	Aliases:     []string{"a"},
	Usage:       "Add a task",
	Description: "Add a task to todoist",
	Arguments: []cli.Argument{
		&cli.StringArg{
			Name:        "content",
			Min:         1,
			Max:         1,
			Destination: &taskAddArgs.Content,
			Config:      cli.StringConfig{TrimSpace: true},
		},
	},
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "description",
			Aliases:  []string{"D"},
			Usage:    "Task description",
			OnlyOnce: true,
			Action: func(_ context.Context, _ *cli.Command, v string) error {
				taskAddArgs.Description = &v
				return nil
			},
		},
		&cli.StringFlag{
			Name:     "project",
			Aliases:  []string{"P"},
			Usage:    "Project ID",
			OnlyOnce: true,
			Action: func(_ context.Context, _ *cli.Command, v string) error {
				taskAddArgs.ProjectID = &v
				return nil
			},
		},
		&cli.StringFlag{
			Name:     "due",
			Aliases:  []string{"d"},
			Usage:    "Due date",
			OnlyOnce: true,
			Action: func(_ context.Context, _ *cli.Command, v string) error {
				taskAddArgs.Due = &sync.Due{String: &v}
				return nil
			},
		},
		&cli.TimestampFlag{
			Name:     "deadline",
			Usage:    "Deadline date",
			OnlyOnce: true,
			Config: cli.TimestampConfig{
				Layouts: []string{time.DateOnly},
			},
			Action: func(_ context.Context, _ *cli.Command, v time.Time) error {
				taskAddArgs.Deadline = &sync.Deadline{Date: v, Lang: cfg.Lang}
				return nil
			},
		},
		&cli.IntFlag{
			Name:     "priority",
			Aliases:  []string{"p"},
			Usage:    "Task priority",
			Value:    1,
			OnlyOnce: true,
			Action: func(_ context.Context, _ *cli.Command, v int64) error {
				taskAddArgs.Priority = utils.IntPtr(int(v))
				return nil
			},
		},
		&cli.StringFlag{
			Name:     "parent",
			Usage:    "Parent task ID",
			OnlyOnce: true,
			Action: func(_ context.Context, _ *cli.Command, v string) error {
				taskAddArgs.ParentID = &v
				return nil
			},
		},
		&cli.StringFlag{
			Name:     "section",
			Aliases:  []string{"s"},
			Usage:    "Section ID",
			OnlyOnce: true,
			Action: func(_ context.Context, _ *cli.Command, v string) error {
				taskAddArgs.SectionID = &v
				return nil
			},
		},
		&cli.StringSliceFlag{
			Name:     "labels",
			Aliases:  []string{"l"},
			Usage:    "Labels",
			OnlyOnce: true,
			Action: func(_ context.Context, _ *cli.Command, v []string) error {
				taskAddArgs.Labels = v
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
		if _, err := cli.Call(ctx, daemon.TaskAdd, taskAddArgs); err != nil {
			return err
		}

		fmt.Printf("Task added: %s\n", taskAddArgs.Content)

		return nil
	},
}
