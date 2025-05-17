package sync

import (
	"fmt"
	"time"

	"github.com/CnTeng/todoist-cli/internal/cmd/util"
	"github.com/CnTeng/todoist-cli/internal/daemon"
	"github.com/spf13/cobra"
)

func NewCmd(f *util.Factory) *cobra.Command {
	params := &daemon.SyncArgs{
		Since: time.Now().AddDate(0, -1, 0),
	}
	cmd := &cobra.Command{
		Use:               "sync",
		Short:             "Sync data",
		Long:              "Sync data with Todoist.",
		Example:           "  todoist sync --all",
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := f.Dial(); err != nil {
				return err
			}
			defer f.Close()

			if _, err := f.Call(cmd.Context(), daemon.Sync, params); err != nil {
				return err
			}

			fmt.Println("Sync success")
			return nil
		},
	}

	cmd.Flags().BoolVarP(&params.All, "all", "a", false, "Sync including completed tasks and archived projects")
	cmd.Flags().BoolVarP(&params.Force, "force", "f", false, "Force sync by removing local data")
	cmd.Flags().AddFlag(newSinceFlag(&params.Since))
	cmd.Flags().BoolP("help", "h", false, "Show help for this command")

	return cmd
}
