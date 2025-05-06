package daemon

import (
	"github.com/CnTeng/todoist-cli/internal/cmd/util"
	"github.com/CnTeng/todoist-cli/internal/daemon"
	"github.com/CnTeng/todoist-cli/internal/db"
	"github.com/spf13/cobra"
)

func NewCmd(f *util.Factory) *cobra.Command {
	return &cobra.Command{
		Use:          "daemon",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
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
}
