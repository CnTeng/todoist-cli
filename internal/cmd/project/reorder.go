package project

import (
	"fmt"
	"slices"

	"github.com/CnTeng/todoist-api-go/sync"
	"github.com/CnTeng/todoist-cli/internal/cmd/util"
	"github.com/CnTeng/todoist-cli/internal/daemon"
	"github.com/CnTeng/todoist-cli/internal/view"
	"github.com/spf13/cobra"
)

func NewReorderCmd(f *util.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "reorder",
		Short:             "Reorder projects",
		Long:              "Reorder projects in Todoist, similar to the 'git rebase -i' command in shell.",
		Example:           "  todoist project reorder",
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := f.Dial(); err != nil {
				return err
			}
			defer f.Close()

			projects := []*sync.Project{}
			if err := f.CallResult(cmd.Context(), daemon.ProjectList, nil, &projects); err != nil {
				return err
			}

			slices.SortFunc(projects, func(i, j *sync.Project) int {
				return i.ChildOrder - j.ChildOrder
			})

			rItems := make([]*view.ReorderItem, 0, len(projects))
			for _, project := range projects {
				// Skip inbox project and archived projects
				if project.InboxProject || project.IsArchived {
					continue
				}

				desc := project.Name
				if project.Description != "" {
					desc = fmt.Sprintf("%s: %s", project.Name, project.Description)
				}

				rItems = append(rItems, &view.ReorderItem{
					ID:          project.ID,
					Description: desc,
				})
			}

			rResult, err := view.NewReorderView(rItems).Reorder()
			if err != nil {
				return err
			}

			params := &sync.ProjectReorderArgs{}
			for i, id := range rResult {
				params.Items = append(params.Items, sync.ProjectReorderItem{
					ID:         id,
					ChildOrder: i,
				})
			}

			if _, err := f.Call(cmd.Context(), daemon.ProjectReorder, params); err != nil {
				return err
			}

			fmt.Println("Projects reordered")
			return nil
		},
	}

	cmd.Flags().BoolP("help", "h", false, "Show help for this command")

	return cmd
}
