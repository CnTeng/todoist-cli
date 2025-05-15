package section

import (
	"fmt"

	"github.com/CnTeng/todoist-api-go/sync"
	"github.com/CnTeng/todoist-cli/internal/cmd/util"
	"github.com/CnTeng/todoist-cli/internal/daemon"
	"github.com/spf13/cobra"
)

func NewAddCmd(f *util.Factory) *cobra.Command {
	params := &sync.SectionAddArgs{}
	cmd := &cobra.Command{
		Use:               "add [flags] --project <project-id> <section>",
		Aliases:           []string{"a"},
		Short:             "Add section",
		Long:              "Add a section to Todoist.",
		Example:           "  todoist section add stage1 --project 6Mjq009rjRw9jXH6",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := f.Dial(); err != nil {
				return err
			}
			defer f.Close()

			params.Name = args[0]
			if _, err := f.Call(cmd.Context(), daemon.SectionAdd, params); err != nil {
				return err
			}

			fmt.Printf("section added: %s\n", params.Name)
			return nil
		},
	}

	cmd.Flags().StringVarP(&params.ProjectID, "project", "P", "", "section project")

	_ = cmd.RegisterFlagCompletionFunc("project", f.NewProjectCompletionFunc(-1))

	_ = cmd.MarkFlagRequired("project")

	return cmd
}
