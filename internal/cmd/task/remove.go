package task

import (
	"fmt"
	"strings"

	"github.com/CnTeng/todoist-api-go/sync"
	"github.com/CnTeng/todoist-cli/internal/cmd/util"
	"github.com/CnTeng/todoist-cli/internal/daemon"
	"github.com/spf13/cobra"
)

func NewRemoveCmd(f *util.Factory) *cobra.Command {
	return &cobra.Command{
		Use:               "remove",
		Aliases:           []string{"rm"},
		Short:             "Remove a task",
		Long:              "Remove a task in todoist",
		Args:              cobra.MinimumNArgs(1),
		ValidArgsFunction: taskCompletion(f),
		RunE: func(cmd *cobra.Command, args []string) error {
			params := make([]*sync.TaskDeleteArgs, 0, len(args))
			for _, arg := range args {
				params = append(params, &sync.TaskDeleteArgs{ID: arg})
			}

			if _, err := f.Call(cmd.Context(), daemon.TaskRemove, params); err != nil {
				return err
			}

			fmt.Printf("Task deleted: %s\n", strings.Join(args, ", "))
			return nil
		},
	}
}
