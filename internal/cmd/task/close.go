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
		Use:        "close",
		Aliases:    []string{"done"},
		Short:      "Close a task",
		Long:       "Close a task in todoist",
		Args:       cobra.MinimumNArgs(1),
		ArgAliases: []string{"id"},

		Run: func(cmd *cobra.Command, args []string) {
			for _, arg := range args {
				if _, err := f.Call(cmd.Context(), daemon.TaskClose, &sync.TaskCloseArgs{ID: arg}); err != nil {
					cobra.CheckErr(err)
				}
				fmt.Printf("Task close: %s\n", arg)
			}
		},
	}
}
