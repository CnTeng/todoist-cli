package task

import (
	"fmt"

	"github.com/CnTeng/todoist-api-go/sync"
	"github.com/CnTeng/todoist-cli/internal/cmd/util"
	"github.com/CnTeng/todoist-cli/internal/daemon"
	"github.com/spf13/cobra"
)

func NewModifyCmd(f *util.Factory, group string) *cobra.Command {
	params := &sync.TaskUpdateArgs{}
	cmd := &cobra.Command{
		Use:               "modify [flags] <task-id>",
		Aliases:           []string{"m"},
		Short:             "Modify task",
		Long:              "Modify a task in Todoist.",
		Example:           "  todoist modify 6X7rM8997g3RQmvh --priority 2",
		GroupID:           group,
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

	contentFlag := newContentFlag(&params.Content)
	deadlineFlag := newDeadlineFlag(&params.Deadline)
	descriptionFlag := newDescriptionFlag(&params.Description)
	dueFlag := newDueFlag(&params.Due)
	durationFlag := newDurationFlag(&params.Duration)
	priorityFlag := newPriorityFlag(&params.Priority)

	cmd.Flags().AddFlag(contentFlag)
	cmd.Flags().AddFlag(deadlineFlag)
	cmd.Flags().AddFlag(descriptionFlag)
	cmd.Flags().AddFlag(dueFlag)
	cmd.Flags().AddFlag(durationFlag)
	cmd.Flags().AddFlag(priorityFlag)
	cmd.Flags().StringSliceVarP(&params.Labels, "labels", "l", nil, "Set the labels of the task (e.g. Food,Shopping)")
	cmd.Flags().BoolP("help", "h", false, "Show help for this command")

	_ = cmd.RegisterFlagCompletionFunc(contentFlag.Name, cobra.NoFileCompletions)
	_ = cmd.RegisterFlagCompletionFunc(deadlineFlag.Name, f.NewDeadlineCompletionFunc(-1))
	_ = cmd.RegisterFlagCompletionFunc(descriptionFlag.Name, cobra.NoFileCompletions)
	_ = cmd.RegisterFlagCompletionFunc(dueFlag.Name, cobra.NoFileCompletions)
	_ = cmd.RegisterFlagCompletionFunc(durationFlag.Name, cobra.NoFileCompletions)
	_ = cmd.RegisterFlagCompletionFunc("labels", f.NewLabelCompletionFunc(-1))
	_ = cmd.RegisterFlagCompletionFunc(priorityFlag.Name, f.NewPriorityCompletionFunc(-1))

	cmd.MarkFlagsOneRequired(contentFlag.Name, deadlineFlag.Name, descriptionFlag.Name, dueFlag.Name, durationFlag.Name, "labels", priorityFlag.Name)

	return cmd
}
