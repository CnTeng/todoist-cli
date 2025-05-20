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
		Use:     "modify [flags] <project-id>",
		Aliases: []string{"m"},
		Short:   "Modify project",
		Long:    "Modify project in Todoist.",
		Example: `  todoist project modify 6Jf8VQXxpwv56VQ7 --name Shopping
  todoist project modify 6Jf8VQXxpwv56VQ7 --favorite=false`,
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: f.NewProjectCompletionFunc(1, nil),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := f.Dial(); err != nil {
				return err
			}
			defer f.Close()

			params.ID = args[0]
			if _, err := f.Call(cmd.Context(), daemon.ProjectModify, params); err != nil {
				return err
			}

			fmt.Printf("Project modified: %s\n", params.ID)
			return nil
		},
	}

	colorFlag := newColorFlag(&params.Color)
	favoriteFlag := newFavoriteFlag(&params.IsFavorite)
	nameFlag := newNameFlag(&params.Name)

	cmd.Flags().AddFlag(colorFlag)
	cmd.Flags().AddFlag(favoriteFlag)
	cmd.Flags().AddFlag(nameFlag)
	cmd.Flags().BoolP("help", "h", false, "Show help for this command")

	_ = cmd.RegisterFlagCompletionFunc(colorFlag.Name, f.NewColorCompletionFunc(-1))
	_ = cmd.RegisterFlagCompletionFunc(favoriteFlag.Name, f.NewFavoriteCompletionFunc(-1))
	_ = cmd.RegisterFlagCompletionFunc(nameFlag.Name, cobra.NoFileCompletions)

	cmd.MarkFlagsOneRequired(colorFlag.Name, favoriteFlag.Name, nameFlag.Name)

	return cmd
}
