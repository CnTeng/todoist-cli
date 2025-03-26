package daemon

import (
	"context"
	"net"

	"github.com/CnTeng/todoist-cli/internal/client"
	"github.com/creachadair/jrpc2"
	"github.com/creachadair/jrpc2/channel"
	"github.com/creachadair/jrpc2/handler"
	"github.com/creachadair/jrpc2/server"
)

type Daemon struct {
	address string
	*client.Client
}

func NewDaemon(address string, c *client.Client) *Daemon {
	return &Daemon{
		address: address,
		Client:  c,
	}
}

func (d *Daemon) Serve(ctx context.Context) error {
	lst, err := net.Listen(jrpc2.Network(d.address))
	if err != nil {
		return err
	}
	defer lst.Close()

	svc := server.Static(handler.Map{
		"listTasks": handler.New(d.ListTasks),
	})

	return server.Loop(ctx, server.NetAccepter(lst, channel.Line), svc, nil)
}
