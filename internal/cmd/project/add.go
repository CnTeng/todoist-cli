package project

import (
	"fmt"

	"github.com/CnTeng/todoist-api-go/sync"
	"github.com/CnTeng/todoist-cli/internal/cmd/util"
	"github.com/CnTeng/todoist-cli/internal/daemon"
	"github.com/spf13/cobra"
)

func NewAddCmd(f *util.Factory) *cobra.Command {
	params := &sync.ProjectAddArgs{}
	cmd := &cobra.Command{
		Use:               "add",
		Aliases:           []string{"a"},
		Short:             "Add a project",
		Long:              "Add a project to todoist",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := f.Dial(); err != nil {
				return err
			}
			defer f.Close()

			params.Name = args[0]
			if _, err := f.Call(cmd.Context(), daemon.ProjectAdd, params); err != nil {
				return err
			}

			fmt.Printf("Project added: %s\n", params.Name)
			return nil
		},
	}

	addColorFlag(f, cmd, &params.Color)
	addFavoriteFlag(f, cmd, &params.IsFavorite)

	return cmd
}
