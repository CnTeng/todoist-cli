package section

import (
	"fmt"
	"strings"

	"github.com/CnTeng/todoist-api-go/sync"
	"github.com/CnTeng/todoist-cli/internal/cmd/util"
	"github.com/CnTeng/todoist-cli/internal/daemon"
	"github.com/spf13/cobra"
)

func NewRemoveCmd(f *util.Factory) *cobra.Command {
	return &cobra.Command{
		Use:     "remove [flags] <section>...",
		Aliases: []string{"rm"},
		Short:   "Remove sections",
		Long:    "Remove sections in Todoist, similar to the 'rm' command in shell.",
		Example: `  todoist section remove 6Xm5HVVRcX00MCjv
  todoist section rm 6Xm5HVVRcX00MCjv 6XxxpwJ00459cJWg`,
		Args:              cobra.MinimumNArgs(1),
		ValidArgsFunction: f.NewSectionCompletionFunc(-1),
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
}
