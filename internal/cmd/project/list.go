package project

import (
	"fmt"

	"github.com/CnTeng/todoist-api-go/sync"
	"github.com/CnTeng/todoist-cli/internal/cmd/util"
	"github.com/CnTeng/todoist-cli/internal/daemon"
	"github.com/CnTeng/todoist-cli/internal/view"
	"github.com/spf13/cobra"
)

func NewListCmd(f *util.Factory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Run: func(cmd *cobra.Command, args []string) {
			result := []*sync.Project{}
			if err := f.CallResult(cmd.Context(), daemon.ProjectList, nil, &result); err != nil {
				cobra.CheckErr(err)
			}

			v := view.NewProjectView(result, f.IconConfig.Icons)
			fmt.Print(v.Render())
		},
	}
}
