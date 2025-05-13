package label

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
		Short:   "List labels",
		Long:    "List labels in todoist, similar to the 'ls' command in shell.",
		Example: `  todoist label list
  todoist label ls`,
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := f.Dial(); err != nil {
				return err
			}
			defer f.Close()

			labels := []*model.Label{}
			if err := f.CallResult(cmd.Context(), daemon.LabelList, nil, &labels); err != nil {
				return err
			}

			v := view.NewLabelView(labels, f.IconConfig.Icons)
			fmt.Print(v.Render())
			return nil
		},
	}
}
