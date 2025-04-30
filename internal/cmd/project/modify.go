package project

import (
	"fmt"

	"github.com/CnTeng/todoist-api-go/sync"
	"github.com/CnTeng/todoist-cli/internal/cmd/util"
	"github.com/CnTeng/todoist-cli/internal/daemon"
	"github.com/spf13/cobra"
)

func NewModifyCmd(f *util.Factory) *cobra.Command {
	params := &sync.ProjectUpdateArgs{}
	cmd := &cobra.Command{
		Use:     "modify",
		Aliases: []string{"m"},
		Short:   "Modify a task",
		Long:    "Modify a task in todoist",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			params.ID = args[0]

			if _, err := f.Call(cmd.Context(), daemon.ProjectModify, params); err != nil {
				return err
			}

			fmt.Printf("Project modified: %s\n", params.ID)
			return nil
		},
	}

	cmd.Flags().AddFlag(newNameFlag(&params.Name))
	addColorFlag(cmd, &params.Color)
	cmd.Flags().AddFlag(newFavoriteFlag(&params.IsFavorite))

	return cmd
}
