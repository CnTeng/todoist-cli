package project

import (
	"fmt"
	"strings"

	"github.com/CnTeng/todoist-api-go/sync"
	"github.com/CnTeng/todoist-cli/internal/cmd/util"
	"github.com/CnTeng/todoist-cli/internal/daemon"
	"github.com/spf13/cobra"
)

func NewRemoveCmd(f *util.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "remove [flags] <project-id>...",
		Aliases:           []string{"rm"},
		Short:             "Remove projects",
		Long:              "Remove projects in Todoist, similar to the 'rm' command in shell.",
		Example:           "  todoist project remove 6X7fphhgwcXVGccJ",
		Args:              cobra.MinimumNArgs(1),
		ValidArgsFunction: f.NewProjectCompletionFunc(-1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := f.Dial(); err != nil {
				return err
			}
			defer f.Close()

			params := []*sync.ProjectDeleteArgs{}
			for _, id := range args {
				params = append(params, &sync.ProjectDeleteArgs{ID: id})
			}

			if _, err := f.Call(cmd.Context(), daemon.ProjectRemove, params); err != nil {
				return err
			}

			fmt.Printf("Projects deleted: %s\n", strings.Join(args, ", "))
			return nil
		},
	}

	cmd.Flags().BoolP("help", "h", false, "Show help for this command")

	return cmd
}
