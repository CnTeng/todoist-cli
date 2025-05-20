package project

import (
	"fmt"

	"github.com/CnTeng/todoist-api-go/sync"
	"github.com/CnTeng/todoist-cli/internal/cmd/util"
	"github.com/CnTeng/todoist-cli/internal/daemon"
	"github.com/CnTeng/todoist-cli/internal/model"
	"github.com/CnTeng/todoist-cli/internal/view"
	"github.com/spf13/cobra"
)

func NewListCmd(f *util.Factory) *cobra.Command {
	params := &model.ProjectListArgs{}
	cmd := &cobra.Command{
		Use:               "list",
		Aliases:           []string{"ls"},
		Short:             "List projects",
		Long:              "List projects in Todoist, similar to the 'ls' command in shell.",
		Example:           "  todoist project list",
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := f.Dial(); err != nil {
				return err
			}
			defer f.Close()

			projects := []*sync.Project{}
			if err := f.CallResult(cmd.Context(), daemon.ProjectList, params, &projects); err != nil {
				return err
			}

			v := view.NewProjectView(projects, f.IconConfig.Icons)
			fmt.Print(v.Render())
			return nil
		},
	}

	cmd.Flags().BoolVarP(&params.Archived, "all", "a", false, "List all projects include archived")
	cmd.Flags().BoolP("help", "h", false, "Show help for this command")

	return cmd
}
