package filter

import (
	"github.com/CnTeng/todoist-cli/internal/cmd/util"
	"github.com/spf13/cobra"
)

func NewCmd(f *util.Factory, group string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "filter",
		Short: "Filter commands",
		Long:  "A set of commands to manage filter in Todoist.",
		Example: `  todoist filter add Important --query 'priority 1'
  todoist filter list
  todoist filter modify 4638879 --name 'Not Important' --query 'priority 4'
  todoist filter remove 4638879
  todoist filter reorder`,
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
