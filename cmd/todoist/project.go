package cmd

import (
	"context"
	"fmt"
	"net"

	"github.com/CnTeng/todoist-api-go/sync"
	tcli "github.com/CnTeng/todoist-cli/internal/cli"
	"github.com/CnTeng/todoist-cli/internal/daemon"
	"github.com/creachadair/jrpc2"
	"github.com/creachadair/jrpc2/channel"
	"github.com/urfave/cli/v3"
)

var ProjectCmd = &cli.Command{
	Name:     "project",
	Usage:    "project commands",
	Commands: []*cli.Command{projectListCmd},
}

var projectListCmd = &cli.Command{
	Name:    "list",
	Aliases: []string{"ls"},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		conn, err := net.Dial("unix", "@todo.sock")
		if err != nil {
			fmt.Printf("Error dialing daemon: %v\n", err)
		}
		defer conn.Close()
		cli := jrpc2.NewClient(channel.Line(conn, conn), nil)

		result := []*sync.Project{}
		if err := cli.CallResult(ctx, daemon.ListProjects, nil, &result); err != nil {
			fmt.Printf("Error calling List Projects: %v\n", err)
		}

		c := tcli.NewCLI(tcli.Nerd)
		c.PrintProjects(result)

		return nil
	},
}
