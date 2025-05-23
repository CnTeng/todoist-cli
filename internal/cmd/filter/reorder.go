package filter

import (
	"fmt"
	"slices"

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

			slices.SortFunc(filters, func(i, j *sync.Filter) int {
				return i.ItemOrder - j.ItemOrder
			})

			rItems := make([]*view.ReorderItem, 0, len(filters))
			for _, filter := range filters {
				rItems = append(rItems, &view.ReorderItem{
					ID:          filter.ID,
					Description: fmt.Sprintf("%s: %s", filter.Name, filter.Query),
				})
			}

			rResult, err := view.NewReorderView(rItems).Reorder()
			if err != nil {
				return err
			}

			params := &sync.FilterReorderArgs{IDOrderMapping: make(map[string]int)}
			for i, id := range rResult {
				params.IDOrderMapping[id] = i
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
