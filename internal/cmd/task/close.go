package task

import (
	"fmt"
	"strings"

	"github.com/CnTeng/todoist-api-go/sync"
	"github.com/CnTeng/todoist-cli/internal/cmd/util"
	"github.com/CnTeng/todoist-cli/internal/daemon"
	"github.com/spf13/cobra"
)

func NewCloseCmd(f *util.Factory, group string) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "close [flags] <task-id>...",
		Short:             "Close task",
		Long:              "Close a task in Todoist.",
		Example:           "  todoist close 6X7rfFVPjhvv84XG",
		GroupID:           group,
		Args:              cobra.MinimumNArgs(1),
		ValidArgsFunction: f.NewTaskCompletionFunc(-1, nil),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := f.Dial(); err != nil {
				return err
			}
			defer f.Close()

			params := []*sync.TaskCloseArgs{}
			for _, arg := range args {
				params = append(params, &sync.TaskCloseArgs{ID: arg})
			}
			if _, err := f.Call(cmd.Context(), daemon.TaskClose, params); err != nil {
				return err
			}

			fmt.Printf("Tasks closed: %s\n", strings.Join(args, ", "))
			return nil
		},
	}

	cmd.Flags().BoolP("help", "h", false, "Show help for this command")

	return cmd
}
