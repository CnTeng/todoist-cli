package label

import (
	"fmt"

	"github.com/CnTeng/todoist-cli/internal/cmd/util"
	"github.com/CnTeng/todoist-cli/internal/daemon"
	"github.com/spf13/cobra"
)

func NewModifyCmd(f *util.Factory) *cobra.Command {
	params := &daemon.LabelUpdateArgs{}
	cmd := &cobra.Command{
		Use:     "modify [flags] <label>",
		Aliases: []string{"m"},
		Short:   "Modify label",
		Long:    "Modify a label in Todoist. For shared labels, only the name can be modified.",
		Example: `  todoist label modify works --name work
  todoist label modify daily -c blue --favorite
  todoist label m daily --favorite=false`,
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: f.NewLabelCompletionFunc(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := f.Dial(); err != nil {
				return err
			}
			defer f.Close()

			params.Name = args[0]
			if _, err := f.Call(cmd.Context(), daemon.LabelModify, params); err != nil {
				return err
			}

			fmt.Printf("Label modified: %s\n", params.Name)
			return nil
		},
	}

	colorFlag := newColorFlag(&params.Args.Color)
	favoriteFlag := newFavoriteFlag(&params.Args.IsFavorite)
	nameFlag := newNameFlag(&params.Args.Name)

	cmd.Flags().AddFlag(colorFlag)
	cmd.Flags().AddFlag(favoriteFlag)
	cmd.Flags().AddFlag(nameFlag)

	_ = cmd.RegisterFlagCompletionFunc(colorFlag.Name, f.NewColorCompletionFunc(-1))
	_ = cmd.RegisterFlagCompletionFunc(favoriteFlag.Name, f.NewFavoriteCompletionFunc(-1))
	_ = cmd.RegisterFlagCompletionFunc(nameFlag.Name, cobra.NoFileCompletions)

	cmd.MarkFlagsOneRequired(colorFlag.Name, favoriteFlag.Name, nameFlag.Name)

	return cmd
}
