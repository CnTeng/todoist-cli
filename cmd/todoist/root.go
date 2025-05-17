package cmd

import (
	"context"

	"github.com/CnTeng/todoist-cli/internal/cmd/daemon"
	"github.com/CnTeng/todoist-cli/internal/cmd/filter"
	"github.com/CnTeng/todoist-cli/internal/cmd/label"
	"github.com/CnTeng/todoist-cli/internal/cmd/project"
	"github.com/CnTeng/todoist-cli/internal/cmd/section"
	"github.com/CnTeng/todoist-cli/internal/cmd/sync"
	"github.com/CnTeng/todoist-cli/internal/cmd/task"
	"github.com/CnTeng/todoist-cli/internal/cmd/util"
	"github.com/spf13/cobra"
)

func newCmd() (*cobra.Command, error) {
	f, err := util.NewFactory()
	if err != nil {
		return nil, err
	}

	cmd := &cobra.Command{
		Use:   "todoist <command> [subcommand] [flags]",
		Short: "A CLI for Todoist",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return f.LoadConfig()
		},
	}

	cmd.PersistentFlags().StringVar(&f.ConfigPath, "config", f.ConfigPath, "config file")

	taskGroup := &cobra.Group{
		ID:    "task",
		Title: "Task commands:",
	}
	otherGroup := &cobra.Group{
		ID:    "other",
		Title: "Resources commands:",
	}

	cmd.AddGroup(taskGroup, otherGroup)

	cmd.AddCommand(
		// Task commands
		task.NewAddCmd(f, taskGroup.ID),
		task.NewQuickAddCmd(f, taskGroup.ID),
		task.NewCloseCmd(f, taskGroup.ID),
		task.NewListCmd(f, taskGroup.ID),
		task.NewModifyCmd(f, taskGroup.ID),
		task.NewMoveCmd(f, taskGroup.ID),
		task.NewRemoveCmd(f, taskGroup.ID),
		task.NewReopenCmd(f, taskGroup.ID),

		// Other resources commands
		project.NewCmd(f, otherGroup.ID),
		section.NewCmd(f, otherGroup.ID),
		label.NewCmd(f, otherGroup.ID),
		filter.NewCmd(f, otherGroup.ID),

		// Advanced commands
		sync.NewCmd(f),
		daemon.NewCmd(f),
	)

	return cmd, nil
}

func Execute() error {
	cmd, err := newCmd()
	if err != nil {
		return err
	}
	return cmd.ExecuteContext(context.Background())
}
