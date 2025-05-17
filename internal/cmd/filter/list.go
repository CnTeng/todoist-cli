package filter

import (
	"fmt"

	"github.com/CnTeng/todoist-api-go/sync"
	"github.com/CnTeng/todoist-cli/internal/cmd/util"
	"github.com/CnTeng/todoist-cli/internal/daemon"
	"github.com/CnTeng/todoist-cli/internal/view"
	"github.com/spf13/cobra"
)

func NewListCmd(f *util.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "list",
		Aliases:           []string{"ls"},
		Short:             "List filters",
		Long:              "List filters in todoist, similar to the 'ls' command in shell.",
		Example:           "  todoist filter list",
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

			v := view.NewFilterView(filters, f.IconConfig.Icons)
			fmt.Print(v.Render())
			return nil
		},
	}

	cmd.Flags().BoolP("help", "h", false, "Show help for this command")

	return cmd
}
