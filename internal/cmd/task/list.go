package task

import (
	"fmt"

	"github.com/CnTeng/todoist-cli/internal/cmd/util"
	"github.com/CnTeng/todoist-cli/internal/daemon"
	"github.com/CnTeng/todoist-cli/internal/model"
	"github.com/CnTeng/todoist-cli/internal/view"
	"github.com/spf13/cobra"
)

func NewListCmd(f *util.Factory) *cobra.Command {
	params := &view.TaskViewConfig{}
	cmd := &cobra.Command{
		Use:               "list [flags]",
		Aliases:           []string{"ls"},
		Short:             "List tasks",
		Long:              "List tasks in Todoist, similar to the 'ls' command in shell.",
		Example:           "  todoist list --all",
		GroupID:           Group.ID,
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := f.Dial(); err != nil {
				return err
			}
			defer f.Close()

			result := []*model.Task{}
			if err := f.CallResult(cmd.Context(), daemon.TaskList, params, &result); err != nil {
				return err
			}

			v := view.NewTaskView(result, f.IconConfig.Icons, params)
			fmt.Print(v.Render())
			return nil
		},
	}

	cmd.Flags().BoolVarP(&params.Completed, "all", "a", false, "List all tasks include completed")
	cmd.Flags().BoolVarP(&params.Description, "description", "d", false, "List tasks with description")
	cmd.Flags().BoolVarP(&params.Tree, "tree", "t", false, "List tasks in tree format")
	cmd.Flags().BoolP("help", "h", false, "Show help for this command")

	return cmd
}
