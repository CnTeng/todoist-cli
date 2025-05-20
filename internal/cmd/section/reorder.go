package section

import (
	"fmt"

	"github.com/CnTeng/todoist-api-go/sync"
	"github.com/CnTeng/todoist-cli/internal/cmd/util"
	"github.com/CnTeng/todoist-cli/internal/daemon"
	"github.com/CnTeng/todoist-cli/internal/model"
	"github.com/CnTeng/todoist-cli/internal/view"
	"github.com/spf13/cobra"
)

func NewReorderCmd(f *util.Factory) *cobra.Command {
	lParams := &model.SectionListArgs{}
	cmd := &cobra.Command{
		Use:               "reorder",
		Short:             "Reorder sections",
		Long:              "Reorder sections in Todoist, similar to the 'git rebase -i' command in shell.",
		Example:           "  todoist section reorder",
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := f.Dial(); err != nil {
				return err
			}
			defer f.Close()

			sections := []*sync.Section{}
			if err := f.CallResult(cmd.Context(), daemon.SectionList, lParams, &sections); err != nil {
				return err
			}

			rItems := make([]*view.ReorderItem, 0, len(sections))
			for _, section := range sections {
				rItems = append(rItems, &view.ReorderItem{
					ID:          section.ID,
					Description: section.Name,
				})
			}

			rResult, err := view.NewReorderView(rItems).Reorder()
			if err != nil {
				return err
			}

			rParams := &sync.SectionReorderArgs{}
			for i, id := range rResult {
				rParams.Items = append(rParams.Items, sync.SectionReorderItem{
					ID:           id,
					SectionOrder: i,
				})
			}

			if _, err := f.Call(cmd.Context(), daemon.SectionReorder, rParams); err != nil {
				return err
			}

			fmt.Println("Sections reordered")
			return nil
		},
	}

	cmd.Flags().StringVarP(&lParams.ProjectID, "project", "p", "", "Filter sections by <project-id>")
	cmd.Flags().BoolP("help", "h", false, "Show help for this command")

	_ = cmd.RegisterFlagCompletionFunc("project", f.NewProjectCompletionFunc(-1))

	_ = cmd.MarkFlagRequired("project")

	return cmd
}
