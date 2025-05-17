package filter

import (
	"fmt"

	"github.com/CnTeng/todoist-api-go/sync"
	"github.com/CnTeng/todoist-cli/internal/cmd/util"
	"github.com/CnTeng/todoist-cli/internal/daemon"
	"github.com/CnTeng/todoist-cli/internal/view"
	"github.com/spf13/cobra"
)

func NewReorderCmd(f *util.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "reorder",
		Short:             "Reorder filters",
		Long:              "Reorder filters in Todoist, similar to the 'git rebase -i' command in shell.",
		Example:           "  todoist filter reorder",
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := f.Dial(); err != nil {
				return err
			}
			defer f.Close()

			filters := []*sync.Filter{}
			if err := f.CallResult(cmd.Context(), daemon.FilterList, nil, &filters); err != nil {
				return err
			}

			params := &sync.FilterReorderArgs{
				IDOrderMapping: make(map[string]int),
			}
			v := view.NewFilterReorderView(filters, f.IconConfig.Icons, params)
			if err := v.Interact(); err != nil {
				return err
			}

			if _, err := f.Call(cmd.Context(), daemon.FilterReorder, params); err != nil {
				return err
			}

			fmt.Println("Filters reordered")
			return nil
		},
	}

	cmd.Flags().BoolP("help", "h", false, "Show help for this command")

	return cmd
}
