package project

import (
	"github.com/CnTeng/todoist-cli/internal/cmd/util"
	"github.com/spf13/cobra"
)

func NewCmd(f *util.Factory, group string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "project",
		Short: "Project commands",
		Long:  "A set of commands to manage projects in Todoist.",
		Example: `  todoist project add 'Shopping List'
  todoist project archive 6X7fphhgwcXVGccJ
  todoist project unarchive 6X7fphhgwcXVGccJ
  todoist project list
  todoist project modify 6Jf8VQXxpwv56VQ7 --name Shopping
  todoist project remove 6X7fphhgwcXVGccJ`,
		GroupID:           group,
		SilenceUsage:      true,
		ValidArgsFunction: cobra.NoFileCompletions,
	}

	cmd.AddCommand(NewAddCmd(f))
	cmd.AddCommand(NewArchiveCmd(f))
	cmd.AddCommand(NewUnarchiveCmd(f))
	cmd.AddCommand(NewListCmd(f))
	cmd.AddCommand(NewModifyCmd(f))
	cmd.AddCommand(NewRemoveCmd(f))

	cmd.Flags().BoolP("help", "h", false, "Show help for this command")

	return cmd
}
