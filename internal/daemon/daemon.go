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
		GetTask:      handler.New(d.db.GetTask),
		ListTasks:    handler.New(d.db.ListTasks),
		AddTask:      handler.New(d.client.AddItem),
		Sync:         handler.NewPos(d.client.Sync, "isForce"),
		ListProjects: handler.New(d.db.ListProjects),
	})

	return server.Loop(ctx, server.NetAccepter(lst, channel.Line), svc, &server.LoopOptions{
		ServerOptions: &jrpc2.ServerOptions{
			Logger: jrpc2.StdLogger(log.New(log.Writer(), "daemon: ", log.Flags())),
		},
	})
}
