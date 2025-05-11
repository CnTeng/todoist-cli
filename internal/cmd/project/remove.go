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
	ids := []string{}
	return &cobra.Command{
		Use:               "remove",
		Aliases:           []string{"rm"},
		Short:             "Remove a project",
		Long:              "Remove a project in todoist",
		Args:              cobra.MinimumNArgs(1),
		ValidArgsFunction: f.NewProjectCompletionFunc(-1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := f.Dial(); err != nil {
				return err
			}
			defer f.Close()

			params := make([]*sync.ProjectDeleteArgs, 0, len(ids))
			for _, id := range args {
				params = append(params, &sync.ProjectDeleteArgs{ID: id})
			}

			if _, err := f.Call(cmd.Context(), daemon.ProjectRemove, params); err != nil {
				return err
			}

			fmt.Printf("Project deleted: %s\n", strings.Join(ids, ", "))
			return nil
		},
	}
}
