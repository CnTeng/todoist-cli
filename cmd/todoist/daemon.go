package cmd

import (
	"context"
	"fmt"

	"github.com/CnTeng/todoist-cli/internal/daemon"
	"github.com/CnTeng/todoist-cli/internal/db"
	"github.com/urfave/cli/v3"
)

func daemonCmd() *cli.Command {
	return &cli.Command{
		Name: "daemon",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			db, err := db.NewDB("test/todoist.db")
			if err != nil {
				fmt.Printf("Error opening database: %v\n", err)
			}

			if err := db.Migrate(); err != nil {
				fmt.Printf("Error migrating database: %v\n", err)
			}

			daemon := daemon.NewDaemon("@todo.sock", cfg.Token, cfg.WSToken, db)

			if err := daemon.Serve(context.Background()); err != nil {
				fmt.Printf("Error starting daemon: %v\n", err)
			}

			return nil
		},
	}
}
