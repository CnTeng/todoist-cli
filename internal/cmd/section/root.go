package section

import (
	"github.com/CnTeng/todoist-cli/internal/cmd/util"
	"github.com/spf13/cobra"
)

func NewCmd(f *util.Factory, group string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "section",
		Short: "Section commands",
		Long:  "A set of commands to manage filter in Todoist.",
		Example: `  todoist section add Groceries --project 6X7FxXvX84jHphx2
  todoist section archive 6X7FxXvX84jHphx2
  todoist section unarchive 6X7FxXvX84jHphx2
  todoist section list
  todoist section modify 6X7FxXvX84jHphx2 --name Supermarket
  todoist section move 6X7FxXvX84jHphx2 --project 9Bw8VQXxpwv56ZY2
  todoist section remove 6X7FxXvX84jHphx2`,
		GroupID:           group,
		SilenceUsage:      true,
		ValidArgsFunction: cobra.NoFileCompletions,
	}

	cmd.AddCommand(NewAddCmd(f))
	cmd.AddCommand(NewArchiveCmd(f))
	cmd.AddCommand(NewUnarchiveCmd(f))
	cmd.AddCommand(NewListCmd(f))
	cmd.AddCommand(NewModifyCmd(f))
	cmd.AddCommand(NewMoveCmd(f))
	cmd.AddCommand(NewRemoveCmd(f))

	cmd.Flags().BoolP("help", "h", false, "Show help for this command")

	return cmd
}
