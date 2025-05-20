package section

import (
	"fmt"

	"github.com/CnTeng/todoist-api-go/sync"
	"github.com/CnTeng/todoist-cli/internal/cmd/util"
	"github.com/CnTeng/todoist-cli/internal/daemon"
	"github.com/spf13/cobra"
)

func NewModifyCmd(f *util.Factory) *cobra.Command {
	params := &sync.SectionUpdateArgs{}
	cmd := &cobra.Command{
		Use:               "modify [flags] <section-id>",
		Aliases:           []string{"m"},
		Short:             "Modify section",
		Long:              "Modify a section in Todoist.",
		Example:           "  todoist section modify 6X7FxXvX84jHphx2 --name Supermarket",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: f.NewSectionCompletionFunc(1, nil),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := f.Dial(); err != nil {
				return err
			}
			defer f.Close()

			params.ID = args[0]
			if _, err := f.Call(cmd.Context(), daemon.SectionModify, params); err != nil {
				return err
			}

			fmt.Printf("Section modified: %s\n", params.ID)
			return nil
		},
	}

	nameFlag := newNameFlag(&params.Name)

	cmd.Flags().AddFlag(nameFlag)
	cmd.Flags().BoolP("help", "h", false, "Show help for this command")

	_ = cmd.RegisterFlagCompletionFunc(nameFlag.Name, cobra.NoFileCompletions)

	_ = cmd.MarkFlagRequired("name")

	return cmd
}
