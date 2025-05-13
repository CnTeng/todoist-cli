package label

import (
	"fmt"

	"github.com/CnTeng/todoist-api-go/sync"
	"github.com/CnTeng/todoist-cli/internal/cmd/util"
	"github.com/CnTeng/todoist-cli/internal/daemon"
	"github.com/CnTeng/todoist-cli/internal/model"
	"github.com/CnTeng/todoist-cli/internal/view"
	"github.com/spf13/cobra"
)

func NewReorderCmd(f *util.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "reorder",
		Short:             "Reorder labels",
		Long:              "Reorder labels in Todoist, similar to the 'git rebase -i' command in shell.",
		Example:           `  todoist label reorder`,
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

			params := &sync.LabelReorderArgs{
				IDOrderMapping: make(map[string]int),
			}
			v := view.NewReorderLabelView(labels, f.IconConfig.Icons, params)
			if err := v.Interact(); err != nil {
				return err
			}

			if _, err := f.Call(cmd.Context(), daemon.LabelReorder, params); err != nil {
				return err
			}

			fmt.Println("Labels reordered")
			return nil
		},
	}

	return cmd
}
