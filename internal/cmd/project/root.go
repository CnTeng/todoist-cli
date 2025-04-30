package project

import (
	"github.com/CnTeng/todoist-cli/internal/cmd/util"
	"github.com/spf13/cobra"
)

func NewCmd(f *util.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "project",
		Aliases:      []string{"proj"},
		Long:         "project commands",
		SilenceUsage: true,
	}

	cmd.AddCommand(NewListCmd(f))
	cmd.AddCommand(NewAddCmd(f))
	cmd.AddCommand(NewModifyCmd(f))
	cmd.AddCommand(NewRemoveCmd(f))
	return cmd
}
