package task

import (
	"fmt"

	"github.com/CnTeng/todoist-api-go/sync"
	"github.com/CnTeng/todoist-cli/internal/cmd/util"
	"github.com/CnTeng/todoist-cli/internal/daemon"
	"github.com/spf13/cobra"
)

func NewMoveCmd(f *util.Factory, group string) *cobra.Command {
	params := &sync.TaskMoveArgs{}
	cmd := &cobra.Command{
		Use:               "move [flags] <task-id>",
		Aliases:           []string{"mv"},
		Short:             "Move task",
		Long:              "Move a task in Todoist.",
		Example:           "  todoist move 6X7rM8997g3RQmvh --parent 6X7rf9x6pv2FGghW",
		GroupID:           group,
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: f.NewTaskCompletionFunc(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := f.Dial(); err != nil {
				return err
			}
			defer f.Close()

			params.ID = args[0]
			if _, err := f.Call(cmd.Context(), daemon.TaskMove, params); err != nil {
				return err
			}

			fmt.Printf("Task moved: %s\n", params.ID)
			return nil
		},
	}

	sectionFlag := newSectionFlag(&params.SectionID)
	parentFlag := newParentFlag(&params.ParentID)
	projectFlag := newProjectFlag(&params.ProjectID)

	cmd.Flags().AddFlag(sectionFlag)
	cmd.Flags().AddFlag(parentFlag)
	cmd.Flags().AddFlag(projectFlag)
	cmd.Flags().BoolP("help", "h", false, "Show help for this command")

	_ = cmd.RegisterFlagCompletionFunc(sectionFlag.Name, f.NewSectionCompletionFunc(-1))
	_ = cmd.RegisterFlagCompletionFunc(parentFlag.Name, f.NewTaskCompletionFunc(-1))
	_ = cmd.RegisterFlagCompletionFunc(projectFlag.Name, f.NewProjectCompletionFunc(-1))

	cmd.MarkFlagsOneRequired(sectionFlag.Name, parentFlag.Name, projectFlag.Name)
	cmd.MarkFlagsMutuallyExclusive(parentFlag.Name, projectFlag.Name)
	cmd.MarkFlagsMutuallyExclusive(parentFlag.Name, sectionFlag.Name)

	return cmd
}
