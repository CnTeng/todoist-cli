package cmd

import (
	"context"
	"fmt"
	"net"

	"github.com/CnTeng/todoist-cli/internal/daemon"
	"github.com/creachadair/jrpc2"
	"github.com/creachadair/jrpc2/channel"
	"github.com/urfave/cli/v3"
)

var syncArgs = &daemon.SyncArgs{}

var syncCmd = &cli.Command{
	Name: "sync",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:        "force",
			Aliases:     []string{"f"},
			Usage:       "force sync",
			OnlyOnce:    true,
			Value:       false,
			Destination: &syncArgs.IsForce,
		},
		&cli.BoolFlag{
			Name:        "all",
			Aliases:     []string{"a"},
			Usage:       "sync all items",
			OnlyOnce:    true,
			Value:       false,
			Destination: &syncArgs.All,
		},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		conn, err := net.Dial("unix", "@todo.sock")
		if err != nil {
			fmt.Printf("Error dialing daemon: %v\n", err)
		}
		defer conn.Close()
		cli := jrpc2.NewClient(channel.Line(conn, conn), nil)

		if _, err := cli.Call(ctx, daemon.Sync, syncArgs); err != nil {
			fmt.Printf("Error calling sync: %v\n", err)
		}

		return nil
	},
}
