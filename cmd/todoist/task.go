package cmd

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/CnTeng/todoist-api-go/sync/v9"
	tcli "github.com/CnTeng/todoist-cli/internal/cli"
	"github.com/CnTeng/todoist-cli/internal/db"
	"github.com/CnTeng/todoist-cli/internal/model"
	"github.com/creachadair/jrpc2"
	"github.com/creachadair/jrpc2/channel"
	"github.com/spf13/cobra"
	"github.com/urfave/cli/v3"
)

var taskListCmd = &cli.Command{
	Name:    "list",
	Aliases: []string{"ls"},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		conn, err := net.Dial("unix", "@todo.sock")
		if err != nil {
			fmt.Printf("Error dialing daemon: %v\n", err)
		}
		defer conn.Close()
		cli := jrpc2.NewClient(channel.Line(conn, conn), nil)

		resp := []*model.Item{}
		if err := cli.CallResult(ctx, "listTasks", nil, &resp); err != nil {
			fmt.Printf("Error calling taskLists: %v\n", err)
		}

		c := tcli.NewCLI(tcli.Nerd)
		c.PrintTasks(resp)

		return nil
	},
}

var taskAddCmd = &cli.Command{
	Name:    "add",
	Aliases: []string{"a"},
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "content",
			Aliases:  []string{"c"},
			Usage:    "Task content",
			Required: true,
			OnlyOnce: true,
		},
		&cli.StringFlag{
			Name:     "description",
			Aliases:  []string{"d"},
			Usage:    "Task description",
			OnlyOnce: true,
		},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		db, err := db.NewDB("test/todoist.db")
		if err != nil {
			fmt.Printf("Error opening database: %v\n", err)
		}

		if err := db.Migrate(); err != nil {
			fmt.Printf("Error migrating database: %v\n", err)
		}

		c := sync.NewClientWithHandler(http.DefaultClient, cfg.Token, db)

		args := &sync.ItemAddArgs{}
		args.Content = cmd.String("content")
		if cmd.IsSet("description") {
			description := cmd.String("description")
			args.Description = &description
		}

		if _, err := c.AddItem(context.Background(), args); err != nil {
			cobra.CheckErr(err)
		}

		return nil
	},
}
