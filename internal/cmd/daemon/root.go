package daemon

import (
	"context"

	"github.com/CnTeng/todoist-cli/internal/cmd/util"
	"github.com/CnTeng/todoist-cli/internal/daemon"
	"github.com/CnTeng/todoist-cli/internal/db"
	"github.com/spf13/cobra"
)

func NewCmd(f *util.Factory) *cobra.Command {
	return &cobra.Command{
		Use:          "daemon",
		SilenceUsage: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return f.ReadConfig()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			db, err := db.NewDB(f.DataFilePath)
			if err != nil {
				return err
			}

			if err := db.Migrate(); err != nil {
				return err
			}

			daemon := daemon.NewDaemon(db, f.DeamonConfig)
			if err := daemon.LoadTokens(); err != nil {
				return err
			}

			return daemon.Serve(context.Background())
		},
	}
}
