package filter

import (
	"fmt"

	"github.com/CnTeng/todoist-api-go/sync"
	"github.com/CnTeng/todoist-cli/internal/cmd/util"
	"github.com/CnTeng/todoist-cli/internal/daemon"
	"github.com/spf13/cobra"
)

func NewAddCmd(f *util.Factory) *cobra.Command {
	params := &sync.FilterAddArgs{}
	cmd := &cobra.Command{
		Use:     "add [flags] --query <query> <filter>",
		Aliases: []string{"a"},
		Short:   "Add filter",
		Long:    "Add a filter to Todoist.",
		Example: `  todoist filter add daily --query 'today | overdue'
  todoist filter a today -q 'today' -c blue --favorite`,
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := f.Dial(); err != nil {
				return err
			}
			defer f.Close()

			params.Name = args[0]
			if _, err := f.Call(cmd.Context(), daemon.FilterAdd, params); err != nil {
				return err
			}

			fmt.Printf("filter added: %s\n", params.Name)
			return nil
		},
	}

	colorFlag := newColorFlag(&params.Color)
	favoriteFlag := newFavoriteFlag(&params.IsFavorite)

	cmd.Flags().AddFlag(colorFlag)
	cmd.Flags().AddFlag(favoriteFlag)
	cmd.Flags().StringVarP(&params.Query, "query", "q", "", "filter query")

	_ = cmd.RegisterFlagCompletionFunc(colorFlag.Name, f.NewColorCompletionFunc(-1))
	_ = cmd.RegisterFlagCompletionFunc(favoriteFlag.Name, f.NewFavoriteCompletionFunc(-1))
	_ = cmd.RegisterFlagCompletionFunc("query", cobra.NoFileCompletions)

	_ = cmd.MarkFlagRequired("query")

	return cmd
}
