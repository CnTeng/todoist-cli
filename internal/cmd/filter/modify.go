package filter

import (
	"fmt"

	"github.com/CnTeng/todoist-api-go/sync"
	"github.com/CnTeng/todoist-cli/internal/cmd/util"
	"github.com/CnTeng/todoist-cli/internal/daemon"
	"github.com/spf13/cobra"
)

func NewModifyCmd(f *util.Factory) *cobra.Command {
	params := &sync.FilterUpdateArgs{}
	cmd := &cobra.Command{
		Use:     "modify [flags] <filter>",
		Aliases: []string{"m"},
		Short:   "Modify filter",
		Long:    "Modify a filter in Todoist.",
		Example: `  todoist filter modify daily1 --name daily
  todoist filter modify daily -c blue --favorite
  todoist filter m daily --favorite=false`,
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: f.NewFilterCompletionFunc(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := f.Dial(); err != nil {
				return err
			}
			defer f.Close()

			params.ID = args[0]
			if _, err := f.Call(cmd.Context(), daemon.FilterModify, params); err != nil {
				return err
			}

			fmt.Printf("Filter modified: %s\n", params.ID)
			return nil
		},
	}

	colorFlag := newColorFlag(&params.Color)
	favoriteFlag := newFavoriteFlag(&params.IsFavorite)
	nameFlag := newNameFlag(&params.Name)
	queryFlag := newQueryFlag(&params.Query)

	cmd.Flags().AddFlag(colorFlag)
	cmd.Flags().AddFlag(favoriteFlag)
	cmd.Flags().AddFlag(nameFlag)
	cmd.Flags().AddFlag(queryFlag)

	_ = cmd.RegisterFlagCompletionFunc(colorFlag.Name, f.NewColorCompletionFunc(-1))
	_ = cmd.RegisterFlagCompletionFunc(favoriteFlag.Name, f.NewFavoriteCompletionFunc(-1))
	_ = cmd.RegisterFlagCompletionFunc(nameFlag.Name, cobra.NoFileCompletions)
	_ = cmd.RegisterFlagCompletionFunc(queryFlag.Name, cobra.NoFileCompletions)

	cmd.MarkFlagsOneRequired(colorFlag.Name, favoriteFlag.Name, nameFlag.Name, queryFlag.Name)

	return cmd
}
