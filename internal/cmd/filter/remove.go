package filter

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
		Use:     "remove [flags] <filter>...",
		Aliases: []string{"rm"},
		Short:   "Remove filters",
		Long:    "Remove filters in Todoist, similar to the 'rm' command in shell.",
		Example: `  todoist filter remove daily
  todoist filter rm work daily`,
		Args:              cobra.MinimumNArgs(1),
		ValidArgsFunction: f.NewFilterCompletionFunc(-1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := f.Dial(); err != nil {
				return err
			}
			defer f.Close()

			params := []*sync.FilterDeleteArgs{}
			for _, arg := range args {
				params = append(params, &sync.FilterDeleteArgs{ID: arg})
			}
			if _, err := f.Call(cmd.Context(), daemon.FilterRemove, params); err != nil {
				return err
			}

			fmt.Printf("Filters deleted: %s\n", strings.Join(args, ", "))
			return nil
		},
	}
}
