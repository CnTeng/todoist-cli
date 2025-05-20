package task

import (
	"fmt"

	"github.com/CnTeng/todoist-api-go/rest"
	"github.com/CnTeng/todoist-api-go/sync"
	"github.com/CnTeng/todoist-cli/internal/cmd/util"
	"github.com/CnTeng/todoist-cli/internal/daemon"
	"github.com/spf13/cobra"
)

func NewAddCmd(f *util.Factory, group string) *cobra.Command {
	params := &sync.TaskAddArgs{}
	cmd := &cobra.Command{
		Use:               "add [flags] <task-name>",
		Aliases:           []string{"a"},
		Short:             "Add task",
		Long:              "Add a task to Todoist.",
		Example:           "  todoist add 'Buy Milk' --project 6Jf8VQXxpwv56VQ7 --labels 'Food,Shopping'",
		GroupID:           group,
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := f.Dial(); err != nil {
				return err
			}
			defer f.Close()

			params.Content = args[0]
			if _, err := f.Call(cmd.Context(), daemon.TaskAdd, params); err != nil {
				return err
			}

			fmt.Printf("Task added: %s\n", params.Content)
			return nil
		},
	}

	deadlineFlag := newDeadlineFlag(&params.Deadline)
	descriptionFlag := newDescriptionFlag(&params.Description)
	dueFlag := newDueFlag(&params.Due)
	durationFlag := newDurationFlag(&params.Duration)
	parentFlag := newParentFlag(&params.ParentID)
	priorityFlag := newPriorityFlag(&params.Priority)
	projectFlag := newProjectFlag(&params.ProjectID)
	sectionFlag := newSectionFlag(&params.SectionID)

	cmd.Flags().AddFlag(deadlineFlag)
	cmd.Flags().AddFlag(descriptionFlag)
	cmd.Flags().AddFlag(dueFlag)
	cmd.Flags().AddFlag(durationFlag)
	cmd.Flags().AddFlag(parentFlag)
	cmd.Flags().AddFlag(priorityFlag)
	cmd.Flags().AddFlag(projectFlag)
	cmd.Flags().AddFlag(sectionFlag)
	cmd.Flags().StringSliceVarP(&params.Labels, "labels", "l", nil, "Set the labels of the task (e.g. 'Food,Shopping')")
	cmd.Flags().BoolP("help", "h", false, "Show help for this command")

	_ = cmd.RegisterFlagCompletionFunc(deadlineFlag.Name, f.NewDeadlineCompletionFunc(-1))
	_ = cmd.RegisterFlagCompletionFunc(descriptionFlag.Name, cobra.NoFileCompletions)
	_ = cmd.RegisterFlagCompletionFunc(dueFlag.Name, cobra.NoFileCompletions)
	_ = cmd.RegisterFlagCompletionFunc(durationFlag.Name, cobra.NoFileCompletions)
	_ = cmd.RegisterFlagCompletionFunc("labels", f.NewLabelCompletionFunc(-1))
	_ = cmd.RegisterFlagCompletionFunc(parentFlag.Name, f.NewTaskCompletionFunc(-1, nil))
	_ = cmd.RegisterFlagCompletionFunc(priorityFlag.Name, f.NewPriorityCompletionFunc(-1))
	_ = cmd.RegisterFlagCompletionFunc(projectFlag.Name, f.NewProjectCompletionFunc(-1, nil))
	_ = cmd.RegisterFlagCompletionFunc(sectionFlag.Name, f.NewSectionCompletionFunc(-1, nil))

	cmd.MarkFlagsMutuallyExclusive(parentFlag.Name, projectFlag.Name)
	cmd.MarkFlagsMutuallyExclusive(parentFlag.Name, sectionFlag.Name)

	return cmd
}

func NewQuickAddCmd(f *util.Factory, group string) *cobra.Command {
	params := &rest.TaskQuickAddRequest{}
	cmd := &cobra.Command{
		Use:               "quick-add <text>",
		Aliases:           []string{"qa"},
		Short:             "Quick add task",
		Long:              "Quick add a task to Todoist.",
		Example:           `  todoist quick-add 'Buy Milk P2 @Food @Shopping'`,
		GroupID:           group,
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := f.Dial(); err != nil {
				return err
			}
			defer f.Close()

			params.Text = args[0]
			if _, err := f.Call(cmd.Context(), daemon.TaskQuickAdd, params); err != nil {
				return err
			}

			fmt.Printf("Task added: %s\n", params.Text)
			return nil
		},
	}

	cmd.Flags().BoolP("help", "h", false, "Show help for this command")

	return cmd
}
