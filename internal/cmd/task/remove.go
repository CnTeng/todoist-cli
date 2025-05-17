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
	cmd := &cobra.Command{
		Use:               "remove [flags] <task-id>...",
		Aliases:           []string{"rm"},
		Short:             "Remove task",
		Long:              "Remove a task in Todoist, similar to the 'rm' command in shell.",
		Example:           "  todoist remove 6X7rfFVPjhvv84XG",
		GroupID:           Group.ID,
		Args:              cobra.MinimumNArgs(1),
		ValidArgsFunction: f.NewTaskCompletionFunc(-1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := f.Dial(); err != nil {
				return err
			}
			defer f.Close()

			params := []*sync.TaskDeleteArgs{}
			for _, arg := range args {
				params = append(params, &sync.TaskDeleteArgs{ID: arg})
			}
			if _, err := f.Call(cmd.Context(), daemon.TaskRemove, params); err != nil {
				return err
			}

			fmt.Printf("Tasks deleted: %s\n", strings.Join(args, ", "))
			return nil
		},
	}

	cmd.Flags().BoolP("help", "h", false, "Show help for this command")

	return cmd
}
