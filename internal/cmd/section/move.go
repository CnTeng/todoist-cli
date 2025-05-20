package section

import (
	"fmt"

	"github.com/CnTeng/todoist-api-go/sync"
	"github.com/CnTeng/todoist-cli/internal/cmd/util"
	"github.com/CnTeng/todoist-cli/internal/daemon"
	"github.com/spf13/cobra"
)

func NewMoveCmd(f *util.Factory) *cobra.Command {
	params := &sync.SectionMoveArgs{}
	cmd := &cobra.Command{
		Use:               "move [flags] <section-id>",
		Aliases:           []string{"mv"},
		Short:             "Move section",
		Long:              "Move a section in Todoist.",
		Example:           "  todoist section move 6X7FxXvX84jHphx2 --project 9Bw8VQXxpwv56ZY2",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: f.NewSectionCompletionFunc(1, nil),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := f.Dial(); err != nil {
				return err
			}
			defer f.Close()

			params.ID = args[0]
			if _, err := f.Call(cmd.Context(), daemon.SectionMove, params); err != nil {
				return err
			}

			fmt.Printf("Section moved: %s\n", params.ID)
			return nil
		},
	}

	projectFlag := newProjectFlag(&params.ProjectID)

	cmd.Flags().AddFlag(projectFlag)
	cmd.Flags().BoolP("help", "h", false, "Show help for this command")

	_ = cmd.RegisterFlagCompletionFunc(projectFlag.Name, f.NewProjectCompletionFunc(-1, nil))

	_ = cmd.MarkFlagRequired(projectFlag.Name)

	return cmd
}
