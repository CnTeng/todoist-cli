package label

import (
	"fmt"

	"github.com/CnTeng/todoist-api-go/sync"
	"github.com/CnTeng/todoist-cli/internal/cmd/util"
	"github.com/CnTeng/todoist-cli/internal/daemon"
	"github.com/spf13/cobra"
)

func NewAddCmd(f *util.Factory) *cobra.Command {
	params := &sync.LabelAddArgs{}
	cmd := &cobra.Command{
		Use:               "add [flags] <label-name>",
		Aliases:           []string{"a"},
		Short:             "Add label",
		Long:              "Add a label to Todoist.",
		Example:           "  todoist label add Food --favorite",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := f.Dial(); err != nil {
				return err
			}
			defer f.Close()

			params.Name = args[0]
			if _, err := f.Call(cmd.Context(), daemon.LabelAdd, params); err != nil {
				return err
			}

			fmt.Printf("Label added: %s\n", params.Name)
			return nil
		},
	}

	colorFlag := newColorFlag(&params.Color)
	favoriteFlag := newFavoriteFlag(&params.IsFavorite)

	cmd.Flags().AddFlag(colorFlag)
	cmd.Flags().AddFlag(favoriteFlag)
	cmd.Flags().BoolP("help", "h", false, "Show help for this command")

	_ = cmd.RegisterFlagCompletionFunc(colorFlag.Name, f.NewColorCompletionFunc(-1))
	_ = cmd.RegisterFlagCompletionFunc(favoriteFlag.Name, f.NewFavoriteCompletionFunc(-1))

	return cmd
}
