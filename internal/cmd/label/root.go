package label

import (
	"github.com/CnTeng/todoist-cli/internal/cmd/util"
	"github.com/spf13/cobra"
)

func NewCmd(f *util.Factory, group string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "label",
		Short: "Label commands",
		Long:  "A set of commands to manage labels in Todoist.",
		Example: `  todoist label add Food --favorite
  todoist label list
  todoist label modify Food --color berry_red
  todoist label remove Food
  todoist label reorder`,
		GroupID:           group,
		SilenceUsage:      true,
		ValidArgsFunction: cobra.NoFileCompletions,
	}

	cmd.AddCommand(NewAddCmd(f))
	cmd.AddCommand(NewListCmd(f))
	cmd.AddCommand(NewModifyCmd(f))
	cmd.AddCommand(NewRemoveCmd(f))
	cmd.AddCommand(NewReorderCmd(f))

	cmd.Flags().BoolP("help", "h", false, "Show help for this command")

	return cmd
}
