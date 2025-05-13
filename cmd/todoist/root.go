package cmd

import (
	"context"

	"github.com/CnTeng/todoist-cli/internal/cmd/daemon"
	"github.com/CnTeng/todoist-cli/internal/cmd/filter"
	"github.com/CnTeng/todoist-cli/internal/cmd/label"
	"github.com/CnTeng/todoist-cli/internal/cmd/project"
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

	cmd.AddGroup(task.Group)

	cmd.AddCommand(
		task.NewAddCmd(f),
		task.NewQuickAddCmd(f),
		task.NewCloseCmd(f),
		task.NewListCmd(f),
		task.NewModifyCmd(f),
		task.NewMoveCmd(f),
		task.NewRemoveCmd(f),
		task.NewReopenCmd(f),
		label.NewCmd(f),
		project.NewCmd(f),
		sync.NewCmd(f),
		daemon.NewCmd(f),
		filter.NewCmd(f),
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
