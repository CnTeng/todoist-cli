package filter

import (
	"github.com/CnTeng/todoist-cli/internal/cmd/util"
	"github.com/spf13/cobra"
)

func NewCmd(f *util.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "filter",
		Short: "Filter commands",
		Long:  "A set of commands to manage filter in Todoist.",
		Example: `  todoist filter add daily --query 'today | overdue'
  todoist filter list
  todoist filter modify daily1 --name daily
  todoist filter remove work daily`,
		SilenceUsage:      true,
		ValidArgsFunction: cobra.NoFileCompletions,
	}

	cmd.AddCommand(NewAddCmd(f))
	cmd.AddCommand(NewListCmd(f))
	cmd.AddCommand(NewModifyCmd(f))
	cmd.AddCommand(NewRemoveCmd(f))
	cmd.AddCommand(NewReorderCmd(f))

	return cmd
}
