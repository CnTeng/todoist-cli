package daemon

import (
	"os"
	"path/filepath"

	"github.com/CnTeng/todoist-cli/internal/cmd/util"
	"github.com/CnTeng/todoist-cli/internal/daemon"
	"github.com/CnTeng/todoist-cli/internal/db"
	"github.com/spf13/cobra"
)

func NewCmd(f *util.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "daemon",
		Short:             "Start daemon",
		Long:              "Start daemon to sync data with Todoist.",
		Example:           "  todoist daemon",
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := os.MkdirAll(filepath.Dir(f.DataPath), 0o755); err != nil {
				return err
			}

			db, err := db.NewDB(f.DataPath)
			if err != nil {
				return err
			}

			if err := db.Migrate(); err != nil {
				return err
			}

			daemon := daemon.NewDaemon(db, f.DeamonConfig)
			return daemon.Serve(cmd.Context())
		},
	}

	cmd.Flags().BoolP("help", "h", false, "Show help for this command")

	return cmd
}
