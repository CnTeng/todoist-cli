package section

import (
	"github.com/CnTeng/todoist-cli/internal/cmd/util"
	"github.com/spf13/cobra"
)

func NewCmd(f *util.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "section",
		Short: "Section commands",
		Long:  "A set of commands to manage filter in Todoist.",
		Example: `  todoist section add stage1 --project 6Mjq009rjRw9jXH6
  todoist section list
	todoist section modify 6Xm5HVVRcX00MCjv --name dailytodoist section modify 6Xm5HVVRcX00MCjv --name stage1
	todoist section remove 6Xm5HVVRcX00MCjv`,
		SilenceUsage:      true,
		ValidArgsFunction: cobra.NoFileCompletions,
	}

	cmd.AddCommand(NewAddCmd(f))
	cmd.AddCommand(NewArchiveCmd(f))
	cmd.AddCommand(NewUnarchiveCmd(f))
	cmd.AddCommand(NewListCmd(f))
	cmd.AddCommand(NewModifyCmd(f))
	cmd.AddCommand(NewRemoveCmd(f))

	return cmd
}
