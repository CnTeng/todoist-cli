package section

import (
	"fmt"
	"strings"

	"github.com/CnTeng/todoist-api-go/sync"
	"github.com/CnTeng/todoist-cli/internal/cmd/util"
	"github.com/CnTeng/todoist-cli/internal/daemon"
	"github.com/CnTeng/todoist-cli/internal/model"
	"github.com/spf13/cobra"
)

func NewRemoveCmd(f *util.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "remove [flags] <section-id>...",
		Aliases:           []string{"rm"},
		Short:             "Remove sections",
		Long:              "Remove sections in Todoist, similar to the 'rm' command in shell.",
		Example:           "  todoist section remove 6X7FxXvX84jHphx2",
		Args:              cobra.MinimumNArgs(1),
		ValidArgsFunction: f.NewSectionCompletionFunc(-1, &model.SectionListArgs{All: true}),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := f.Dial(); err != nil {
				return err
			}
			defer f.Close()

			params := []*sync.SectionDeleteArgs{}
			for _, arg := range args {
				params = append(params, &sync.SectionDeleteArgs{ID: arg})
			}
			if _, err := f.Call(cmd.Context(), daemon.SectionRemove, params); err != nil {
				return err
			}

			fmt.Printf("Sections deleted: %s\n", strings.Join(args, ", "))
			return nil
		},
	}

	cmd.Flags().BoolP("help", "h", false, "Show help for this command")

	return cmd
}
