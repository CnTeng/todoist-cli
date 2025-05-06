package cmd

import (
	"context"

	"github.com/CnTeng/todoist-cli/internal/cmd/daemon"
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
			if err := f.LoadConfig(); err != nil {
				return err
			}

			if err := f.Dial(); err != nil {
				return err
			}

			return nil
		},

		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			return f.Close()
		},
	}

	cmd.PersistentFlags().StringVar(&f.ConfigPath, "config", f.ConfigPath, "config file")

	cmd.AddCommand(
		task.NewListCmd(f),
		task.NewAddCmd(f),
		task.NewQuickAddCmd(f),
		task.NewModifyCmd(f),
		task.NewCloseCmd(f),
		task.NewReopenCmd(f),
		task.NewRemoveCmd(f),
		task.NewMoveCmd(f),
		project.NewCmd(f),
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
