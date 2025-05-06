package task

import (
	"fmt"

	"github.com/CnTeng/todoist-api-go/rest"
	"github.com/CnTeng/todoist-api-go/sync"
	"github.com/CnTeng/todoist-cli/internal/cmd/util"
	"github.com/CnTeng/todoist-cli/internal/daemon"
	"github.com/spf13/cobra"
)

func NewAddCmd(f *util.Factory) *cobra.Command {
	params := &sync.TaskAddArgs{}
	cmd := &cobra.Command{
		Use:     "add",
		Aliases: []string{"a"},
		Short:   "Add a task",
		Long:    "Add a task to todoist",
		GroupID: Group.ID,
		Args:    cobra.ExactArgs(1),
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

	cmd.Flags().AddFlag(newDescriptionFlag(&params.Description))
	cmd.Flags().AddFlag(newProjectFlag(&params.ProjectID))
	cmd.Flags().AddFlag(newDueFlag(&params.Due))
	cmd.Flags().AddFlag(newDeadlineFlag(&params.Deadline))
	cmd.Flags().AddFlag(newPriorityFlag(&params.Priority))
	cmd.Flags().AddFlag(newParentFlag(&params.ParentID))
	cmd.Flags().AddFlag(newSectionFlag(&params.SectionID))
	addLabelsFlag(cmd, &params.Labels)

	return cmd
}

func NewQuickAddCmd(f *util.Factory) *cobra.Command {
	params := &rest.TaskQuickAddRequest{}
	return &cobra.Command{
		Use:     "quick-add",
		Aliases: []string{"qa"},
		Short:   "Quick add a task",
		Long:    "Quick Add a task to todoist",
		GroupID: Group.ID,
		Args:    cobra.ExactArgs(1),
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
}
