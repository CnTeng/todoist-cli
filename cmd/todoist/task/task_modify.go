package task

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/CnTeng/todoist-api-go/sync"
	"github.com/CnTeng/todoist-cli/internal/daemon"
	"github.com/CnTeng/todoist-cli/internal/model"
	"github.com/CnTeng/todoist-cli/internal/utils"
	"github.com/creachadair/jrpc2"
	"github.com/creachadair/jrpc2/channel"
	"github.com/urfave/cli/v3"
)

func NewModifyCmd(cfg *model.Config) *cli.Command {
	params := &sync.ItemUpdateArgs{}
	return &cli.Command{
		Name:        "modify",
		Aliases:     []string{"m"},
		Usage:       "Modify a task",
		Description: "Modify a task in todoist",
		Arguments: []cli.Argument{
			&cli.StringArg{
				Name:        "ID",
				Min:         1,
				Max:         1,
				Destination: &params.ID,
				Config:      cli.StringConfig{TrimSpace: true},
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "content",
				Aliases:  []string{"c"},
				Usage:    "Task content",
				OnlyOnce: true,
				Action: func(_ context.Context, _ *cli.Command, v string) error {
					params.Content = &v
					return nil
				},
			},
			&cli.StringFlag{
				Name:     "description",
				Aliases:  []string{"D"},
				Usage:    "Task description",
				OnlyOnce: true,
				Action: func(_ context.Context, _ *cli.Command, v string) error {
					params.Description = &v
					return nil
				},
			},
			&cli.StringFlag{
				Name:     "due",
				Aliases:  []string{"d"},
				Usage:    "Due date",
				OnlyOnce: true,
				Action: func(_ context.Context, _ *cli.Command, v string) error {
					params.Due = &sync.Due{String: &v}
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
					params.Deadline = &sync.Deadline{Date: v, Lang: cfg.Lang}
					return nil
				},
			},
			&cli.IntFlag{
				Name:     "priority",
				Aliases:  []string{"p"},
				Usage:    "Task priority",
				OnlyOnce: true,
				Action: func(_ context.Context, _ *cli.Command, v int64) error {
					params.Priority = utils.IntPtr(int(v))
					return nil
				},
			},
			&cli.StringSliceFlag{
				Name:     "labels",
				Aliases:  []string{"l"},
				Usage:    "Labels",
				OnlyOnce: true,
				Action: func(_ context.Context, _ *cli.Command, v []string) error {
					params.Labels = v
					return nil
				},
			},
			&cli.StringFlag{
				Name:     "duration",
				Usage:    "Duration",
				OnlyOnce: true,
				Action: func(_ context.Context, _ *cli.Command, v string) error {
					duration, err := sync.ParseDuration(v)
					if err != nil {
						fmt.Printf("Error parsing duration: %v\n", err)
					}
					params.Duration = duration
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
			if _, err := cli.Call(ctx, daemon.TaskModify, params); err != nil {
				return err
			}

			fmt.Printf("Task modified: %s\n", params.ID)

			return nil
		},
	}
}
