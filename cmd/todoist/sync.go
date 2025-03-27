package cmd

import (
	"context"
	"fmt"
	"net"

	"github.com/creachadair/jrpc2"
	"github.com/creachadair/jrpc2/channel"
	"github.com/urfave/cli/v3"
)

var syncCmd = &cli.Command{
	Name: "sync",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		conn, err := net.Dial("unix", "@todo.sock")
		if err != nil {
			fmt.Printf("Error dialing daemon: %v\n", err)
		}
		defer conn.Close()
		cli := jrpc2.NewClient(channel.Line(conn, conn), nil)

		if _, err := cli.Call(ctx, "sync", nil); err != nil {
			fmt.Printf("Error calling taskLists: %v\n", err)
		}

		return nil
	},
}
