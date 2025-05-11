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
		Use:               "add",
		Aliases:           []string{"a"},
		Short:             "Add a task",
		Long:              "Add a task to todoist",
		GroupID:           Group.ID,
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

	addDeadlineFlag(f, cmd, &params.Deadline)
	addDescriptionFlag(f, cmd, &params.Description)
	addDueFlag(f, cmd, &params.Due)
	addLabelsFlag(f, cmd, &params.Labels)
	addParentFlag(f, cmd, &params.ParentID)
	addPriorityFlag(f, cmd, &params.Priority)
	addProjectFlag(f, cmd, &params.ProjectID)
	addSectionFlag(f, cmd, &params.SectionID)

	return cmd
}

func NewQuickAddCmd(f *util.Factory) *cobra.Command {
	params := &rest.TaskQuickAddRequest{}
	return &cobra.Command{
		Use:               "quick-add",
		Aliases:           []string{"qa"},
		Short:             "Quick add a task",
		Long:              "Quick Add a task to todoist",
		GroupID:           Group.ID,
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
}
