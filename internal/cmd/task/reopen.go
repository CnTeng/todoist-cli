package task

import (
	"fmt"

	"github.com/CnTeng/todoist-api-go/sync"
	"github.com/CnTeng/todoist-cli/internal/cmd/util"
	"github.com/CnTeng/todoist-cli/internal/daemon"
	"github.com/spf13/cobra"
)

func NewReopenCmd(f *util.Factory) *cobra.Command {
	return &cobra.Command{
		Use:     "reopen",
		Aliases: []string{"r"},
		Short:   "Reopen a task",
		Long:    "Reopen a task in todoist",
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := f.Dial(); err != nil {
				return err
			}
			defer f.Close()

			for _, arg := range args {
				if _, err := f.Call(cmd.Context(), daemon.TaskReopen, &sync.TaskUncompleteArgs{ID: arg}); err != nil {
					return err
				}
				fmt.Printf("Task reopen: %s\n", arg)
			}
			return nil
		},
	}
}
