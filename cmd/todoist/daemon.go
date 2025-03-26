package cmd

import (
	"context"
	"fmt"
	"net/http"

	"github.com/CnTeng/todoist-api-go/sync/v9"
	"github.com/CnTeng/todoist-cli/internal/client"
	"github.com/CnTeng/todoist-cli/internal/daemon"
	"github.com/CnTeng/todoist-cli/internal/db"
	"github.com/urfave/cli/v3"
)

var daemonCmd = &cli.Command{
	Name: "daemon",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		db, err := db.NewDB("test/todoist.db")
		if err != nil {
			fmt.Printf("Error opening database: %v\n", err)
		}

		if err := db.Migrate(); err != nil {
			fmt.Printf("Error migrating database: %v\n", err)
		}

		sc := sync.NewClientWithHandler(http.DefaultClient, cfg.Token, db)
		c := client.NewClient(db, sc)
		daemon := daemon.NewDaemon("@todo.sock", c)

		if err := daemon.Serve(context.Background()); err != nil {
			fmt.Printf("Error starting daemon: %v\n", err)
		}

		return nil
	},
}
