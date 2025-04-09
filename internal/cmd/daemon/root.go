package daemon

import (
	"context"

	"github.com/CnTeng/todoist-cli/internal/cmd/util"
	"github.com/CnTeng/todoist-cli/internal/daemon"
	"github.com/CnTeng/todoist-cli/internal/db"
	"github.com/urfave/cli/v3"
)

func NewCmd(f *util.Factory) *cli.Command {
	return &cli.Command{
		Name: "daemon",
		Action: func(ctx context.Context, cmd *cli.Command) error {
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
