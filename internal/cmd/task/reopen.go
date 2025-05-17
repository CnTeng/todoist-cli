package task

import (
	"fmt"
	"strings"

	"github.com/CnTeng/todoist-api-go/sync"
	"github.com/CnTeng/todoist-cli/internal/cmd/util"
	"github.com/CnTeng/todoist-cli/internal/daemon"
	"github.com/spf13/cobra"
)

func NewReopenCmd(f *util.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "reopen [flags] <task-id>...",
		Short:             "Reopen task",
		Long:              "Reopen a task in Todoist.",
		Example:           "  todoist reopen 6X7rfFVPjhvv84XG",
		GroupID:           Group.ID,
		Args:              cobra.MinimumNArgs(1),
		ValidArgsFunction: f.NewTaskCompletionFunc(-1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := f.Dial(); err != nil {
				return err
			}
			defer f.Close()

			params := []*sync.TaskUncompleteArgs{}
			for _, arg := range args {
				params = append(params, &sync.TaskUncompleteArgs{ID: arg})
			}
			if _, err := f.Call(cmd.Context(), daemon.TaskReopen, params); err != nil {
				return err
			}

			fmt.Printf("Task reopen: %s\n", strings.Join(args, ", "))
			return nil
		},
	}

	cmd.Flags().BoolP("help", "h", false, "Show help for this command")

	return cmd
}
