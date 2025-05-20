package task

import (
	"fmt"

	"github.com/CnTeng/todoist-api-go/sync"
	"github.com/CnTeng/todoist-cli/internal/cmd/util"
	"github.com/CnTeng/todoist-cli/internal/daemon"
	"github.com/CnTeng/todoist-cli/internal/model"
	"github.com/CnTeng/todoist-cli/internal/view"
	"github.com/spf13/cobra"
)

func NewReorderCmd(f *util.Factory, group string) *cobra.Command {
	params := &model.TaskListArgs{}
	cmd := &cobra.Command{
		Use:               "reorder",
		Short:             "Reorder tasks",
		Long:              "Reorder tasks in Todoist, similar to the 'git rebase -i' command in shell.",
		Example:           "  todoist reorder",
		GroupID:           group,
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := f.Dial(); err != nil {
				return err
			}
			defer f.Close()

			tasks := []*model.Task{}
			if err := f.CallResult(cmd.Context(), daemon.TaskList, params, &tasks); err != nil {
				return err
			}

			rItems := make([]*view.ReorderItem, 0, len(tasks))
			for _, task := range tasks {
				// Skip checked tasks
				if task.Checked {
					continue
				}

				rItems = append(rItems, &view.ReorderItem{
					ID:          task.ID,
					Description: task.Content,
				})
			}

			rResult, err := view.NewReorderView(rItems).Reorder()
			if err != nil {
				return err
			}

			params := &sync.TaskReorderArgs{}
			for i, id := range rResult {
				params.Items = append(params.Items, sync.TaskReorderItem{
					ID:         id,
					ChildOrder: i,
				})
			}

			if _, err := f.Call(cmd.Context(), daemon.TaskReorder, params); err != nil {
				return err
			}

			fmt.Println("Tasks reordered")
			return nil
		},
	}

	cmd.Flags().StringVarP(&params.ProjectID, "project", "p", "", "Filter tasks by <project-id>")
	cmd.Flags().BoolP("help", "h", false, "Show help for this command")

	_ = cmd.RegisterFlagCompletionFunc("project", f.NewProjectCompletionFunc(-1))

	_ = cmd.MarkFlagRequired("project")

	return cmd
}
