package task

import (
	"context"
	"fmt"

	"github.com/CnTeng/todoist-api-go/rest"
	"github.com/CnTeng/todoist-api-go/sync"
	"github.com/CnTeng/todoist-cli/internal/cmd/util"
	"github.com/CnTeng/todoist-cli/internal/daemon"
	"github.com/urfave/cli/v3"
)

func NewAddCmd(f *util.Factory) *cli.Command {
	params := &sync.TaskAddArgs{}
	return &cli.Command{
		Name:        "add",
		Aliases:     []string{"a"},
		Usage:       "Add a task",
		Description: "Add a task to todoist",
		Category:    "task",
		Arguments: []cli.Argument{
			&cli.StringArg{
				Name:        "content",
				Destination: &params.Content,
				Config:      cli.StringConfig{TrimSpace: true},
			},
		},
		Flags: []cli.Flag{
			newDescriptionFlag(&params.Description),
			newProjectFlag(&params.ProjectID),
			newDueFlag(&params.Due),
			newDeadlineFlag(&params.Deadline),
			newPriorityFlag(&params.Priority),
			newParentFlag(&params.ParentID),
			newSectionFlag(&params.SectionID),
			newLabelsFlag(&params.Labels),
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if _, err := f.Call(ctx, daemon.TaskAdd, params); err != nil {
				return err
			}

			fmt.Printf("Task added: %s\n", params.Content)
			return nil
		},
	}
}

func NewQuickAddCmd(f *util.Factory) *cli.Command {
	params := &rest.TaskQuickAddRequest{}
	return &cli.Command{
		Name:        "quick-add",
		Aliases:     []string{"qa"},
		Usage:       "Quick add a task",
		Description: "Quick Add a task to todoist",
		Category:    "task",
		Arguments: []cli.Argument{
			&cli.StringArg{
				Name:        "text",
				Destination: &params.Text,
				Config:      cli.StringConfig{TrimSpace: true},
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if _, err := f.Call(ctx, daemon.TaskQuickAdd, params); err != nil {
				return err
			}

			fmt.Printf("Task added: %s\n", params.Text)
			return nil
		},
	}
}
