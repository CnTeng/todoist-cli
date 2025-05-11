package task

import (
	"fmt"

	"github.com/CnTeng/todoist-api-go/sync"
	"github.com/CnTeng/todoist-cli/internal/cmd/util"
	"github.com/CnTeng/todoist-cli/internal/daemon"
	"github.com/spf13/cobra"
)

func NewCloseCmd(f *util.Factory) *cobra.Command {
	return &cobra.Command{
		Use:               "close",
		Aliases:           []string{"done"},
		Short:             "Close a task",
		Long:              "Close a task in todoist",
		GroupID:           Group.ID,
		Args:              cobra.MinimumNArgs(1),
		ArgAliases:        []string{"id"},
		ValidArgsFunction: f.NewTaskCompletionFunc(-1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := f.Dial(); err != nil {
				return err
			}
			defer f.Close()

			for _, arg := range args {
				if _, err := f.Call(cmd.Context(), daemon.TaskClose, &sync.TaskCloseArgs{ID: arg}); err != nil {
					return err
				}
				fmt.Printf("Task close: %s\n", arg)
			}
			return nil
		},
	}
}
