package task

import (
	"fmt"

	"github.com/CnTeng/todoist-api-go/sync"
	"github.com/CnTeng/todoist-cli/internal/cmd/util"
	"github.com/CnTeng/todoist-cli/internal/daemon"
	"github.com/spf13/cobra"
)

func NewModifyCmd(f *util.Factory) *cobra.Command {
	params := &sync.TaskUpdateArgs{}
	cmd := &cobra.Command{
		Use:               "modify",
		Aliases:           []string{"m"},
		Short:             "Modify a task",
		Long:              "Modify a task in todoist",
		GroupID:           Group.ID,
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: f.NewTaskCompletionFunc(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := f.Dial(); err != nil {
				return err
			}
			defer f.Close()

			params.ID = args[0]
			if _, err := f.Call(cmd.Context(), daemon.TaskModify, params); err != nil {
				return err
			}

			fmt.Printf("Task modified: %s\n", params.ID)
			return nil
		},
	}

	addContentFlag(f, cmd, &params.Content)
	addDeadlineFlag(f, cmd, &params.Deadline)
	addDescriptionFlag(f, cmd, &params.Description)
	addDueFlag(f, cmd, &params.Due)
	addDurationFlag(f, cmd, &params.Duration)
	addLabelsFlag(f, cmd, &params.Labels)
	addPriorityFlag(f, cmd, &params.Priority)

	return cmd
}
