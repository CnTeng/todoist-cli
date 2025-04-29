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
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List tasks",
		Long:    "List tasks in todoist",
		Run: func(cmd *cobra.Command, args []string) {
			result := []*model.Task{}
			if err := f.CallResult(cmd.Context(), daemon.TaskList, params, &result); err != nil {
				cobra.CheckErr(err)
			}

			v := view.NewTaskView(result, f.IconConfig.Icons, params)
			fmt.Print(v.Render())
		},
	}

	cmd.Flags().BoolVarP(&params.Completed, "all", "a", true, "list all tasks include completed")
	cmd.Flags().BoolVarP(&params.Description, "description", "d", true, "list tasks include description")
	cmd.Flags().BoolVarP(&params.Tree, "tree", "t", true, "list tasks in tree format")

	return cmd
}
