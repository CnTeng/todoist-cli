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
		Use:     "add",
		Aliases: []string{"a"},
		Short:   "Add a project",
		Long:    "Add a project to todoist",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			params.Name = args[0]
			if _, err := f.Call(cmd.Context(), daemon.ProjectAdd, params); err != nil {
				return err
			}

			fmt.Printf("Project added: %s\n", params.Name)
			return nil
		},
	}

	addColorFlag(cmd, &params.Color)
	cmd.Flags().AddFlag(newFavoriteFlag(&params.IsFavorite))

	return cmd
}
