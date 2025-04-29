package task

import (
	"fmt"

	"github.com/CnTeng/todoist-api-go/sync"
	"github.com/CnTeng/todoist-cli/internal/cmd/util"
	"github.com/CnTeng/todoist-cli/internal/daemon"
	"github.com/spf13/cobra"
)

func NewMoveCmd(f *util.Factory) *cobra.Command {
	params := &sync.TaskMoveArgs{}
	cmd := &cobra.Command{
		Use:     "move",
		Aliases: []string{"mv"},
		Short:   "Move a task",
		Long:    "Move a task in todoist",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			params.ID = args[0]
			if _, err := f.Call(cmd.Context(), daemon.TaskMove, params); err != nil {
				cobra.CheckErr(err)
			}

			fmt.Printf("Task moved: %s\n", params.ID)
		},
	}

	cmd.Flags().AddFlag(newSectionFlag(&params.SectionID))
	cmd.Flags().AddFlag(newParentFlag(&params.ParentID))
	cmd.Flags().AddFlag(newProjectFlag(&params.ProjectID))

	return cmd
}
