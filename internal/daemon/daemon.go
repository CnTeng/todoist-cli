package daemon

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/CnTeng/todoist-api-go/sync/v9"
	"github.com/CnTeng/todoist-cli/internal/db"
	"github.com/creachadair/jrpc2"
	"github.com/creachadair/jrpc2/channel"
	"github.com/creachadair/jrpc2/handler"
	"github.com/creachadair/jrpc2/server"
)

type Daemon struct {
	address string
	client  *sync.Client
	db      *db.DB
}

func NewDaemon(address string, token string, db *db.DB) *Daemon {
	return &Daemon{
		address: address,
		client:  sync.NewClientWithHandler(http.DefaultClient, token, db),
		db:      db,
	}
}

func (d *Daemon) Serve(ctx context.Context) error {
	lst, err := net.Listen(jrpc2.Network(d.address))
	if err != nil {
		return err
	}
	defer lst.Close()

	svc := server.Static(handler.Map{
		"getTask":   handler.New(d.db.GetTask),
		"listTasks": handler.New(d.db.ListTasks),
		"sync":      handler.New(d.client.Sync),
	})

	return server.Loop(ctx, server.NetAccepter(lst, channel.Line), svc, &server.LoopOptions{
		ServerOptions: &jrpc2.ServerOptions{
			Logger: jrpc2.StdLogger(log.New(log.Writer(), "daemon: ", log.Flags())),
		},
	})
}
