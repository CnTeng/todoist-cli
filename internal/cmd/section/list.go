package section

import (
	"fmt"

	"github.com/CnTeng/todoist-cli/internal/cmd/util"
	"github.com/CnTeng/todoist-cli/internal/daemon"
	"github.com/CnTeng/todoist-cli/internal/model"
	"github.com/CnTeng/todoist-cli/internal/view"
	"github.com/spf13/cobra"
)

func NewListCmd(f *util.Factory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List sections",
		Long:    "List sections in todoist, similar to the 'ls' command in shell.",
		Example: `  todoist section list
  todoist section ls`,
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := f.Dial(); err != nil {
				return err
			}
			defer f.Close()

			sections := []*model.Section{}
			if err := f.CallResult(cmd.Context(), daemon.SectionList, nil, &sections); err != nil {
				return err
			}

			v := view.NewSectionView(sections, f.IconConfig.Icons)
			fmt.Print(v.Render())
			return nil
		},
	}
}
