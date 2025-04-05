package task

import (
	"context"
	"fmt"
	"time"

	"github.com/CnTeng/todoist-api-go/sync"
	"github.com/CnTeng/todoist-cli/cmd/todoist/util"
	"github.com/CnTeng/todoist-cli/internal/daemon"
	"github.com/CnTeng/todoist-cli/internal/utils"
	"github.com/urfave/cli/v3"
)

func NewAddCmd(f *util.Factory) *cli.Command {
	params := &sync.ItemAddArgs{}
	return &cli.Command{
		Name:        "add",
		Aliases:     []string{"a"},
		Usage:       "Add a task",
		Description: "Add a task to todoist",
		Category:    "task",
		Arguments: []cli.Argument{
			&cli.StringArg{
				Name:        "content",
				Min:         1,
				Max:         1,
				Destination: &params.Content,
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
					params.Description = &v
					return nil
				},
			},
			&cli.StringFlag{
				Name:     "project",
				Aliases:  []string{"P"},
				Usage:    "Project ID",
				OnlyOnce: true,
				Action: func(_ context.Context, _ *cli.Command, v string) error {
					params.ProjectID = &v
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
					params.Deadline = &sync.Deadline{Date: v, Lang: f.Config.Lang}
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
					params.Priority = utils.IntPtr(int(v))
					return nil
				},
			},
			&cli.StringFlag{
				Name:     "parent",
				Usage:    "Parent task ID",
				OnlyOnce: true,
				Action: func(_ context.Context, _ *cli.Command, v string) error {
					params.ParentID = &v
					return nil
				},
			},
			&cli.StringFlag{
				Name:     "section",
				Aliases:  []string{"s"},
				Usage:    "Section ID",
				OnlyOnce: true,
				Action: func(_ context.Context, _ *cli.Command, v string) error {
					params.SectionID = &v
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
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if _, err := f.RpcClient.Call(ctx, daemon.TaskAdd, params); err != nil {
				return err
			}

			fmt.Printf("Task added: %s\n", params.Content)
			return nil
		},
	}
}

func NewQuickAddCmd(f *util.Factory) *cli.Command {
	params := &sync.ItemQuickAddParams{}
	return &cli.Command{
		Name:        "quick-add",
		Aliases:     []string{"qa"},
		Usage:       "Quick add a task",
		Description: "Quick Add a task to todoist",
		Category:    "task",
		Arguments: []cli.Argument{
			&cli.StringArg{
				Name:        "text",
				Min:         1,
				Max:         1,
				Destination: &params.Text,
				Config:      cli.StringConfig{TrimSpace: true},
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if _, err := f.RpcClient.Call(ctx, daemon.TaskQuickAdd, params); err != nil {
				return err
			}

			fmt.Printf("Task added: %s\n", params.Text)
			return nil
		},
	}
}
