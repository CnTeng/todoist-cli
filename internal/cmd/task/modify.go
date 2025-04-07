package task

import (
	"context"
	"fmt"

	"github.com/CnTeng/todoist-api-go/sync"
	"github.com/CnTeng/todoist-cli/internal/cmd/util"
	"github.com/CnTeng/todoist-cli/internal/daemon"
	"github.com/urfave/cli/v3"
)

func NewModifyCmd(f *util.Factory) *cli.Command {
	params := &sync.ItemUpdateArgs{}
	return &cli.Command{
		Name:        "modify",
		Aliases:     []string{"m"},
		Usage:       "Modify a task",
		Description: "Modify a task in todoist",
		Category:    "task",
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
			newContentFlag(params.Content),
			newDescriptionFlag(params.Description),
			newDueFlag(params.Due),
			newDeadlineFlag(params.Deadline, f.Lang),
			newPriorityFlag(params.Priority),
			newLabelsFlag(&params.Labels),
			newDurationFlag(params.Duration),
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if _, err := f.Call(ctx, daemon.TaskModify, params); err != nil {
				return err
			}

			fmt.Printf("Task modified: %s\n", params.ID)
			return nil
		},
	}
}
