package daemon

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/CnTeng/todoist-api-go/sync"
	"github.com/CnTeng/todoist-api-go/ws"
	"github.com/CnTeng/todoist-cli/internal/db"
	"github.com/creachadair/jrpc2"
	"github.com/creachadair/jrpc2/channel"
	"github.com/creachadair/jrpc2/handler"
	"github.com/creachadair/jrpc2/server"
)

type Daemon struct {
	address string
	client  *sync.Client
	ws      *ws.Client
	db      *db.DB
	log     *log.Logger
}

func NewDaemon(address string, token string, wsToken string, db *db.DB) *Daemon {
	d := &Daemon{
		address: address,
		client:  sync.NewClient(http.DefaultClient, token, db),
		db:      db,
		log:     log.New(log.Writer(), "daemon: ", log.Flags()),
	}
	d.ws = ws.NewClient(wsToken, d)

	return d
}

func (d *Daemon) HandleNotification(ctx context.Context, noti ws.Notification) error {
	if noti != ws.SyncNeeded {
		return nil
	}
	d.log.Println("sync needed")
	_, err := d.client.Sync(ctx, false)
	return err
}

func (d *Daemon) Serve(ctx context.Context) error {
	if err := d.ws.Listen(ctx); err != nil {
		return err
	}

	lst, err := net.Listen(jrpc2.Network(d.address))
	if err != nil {
		return err
	}
	defer lst.Close()

	svc := server.Static(handler.Map{
		GetTask:      handler.New(d.db.GetTask),
		ListTasks:    handler.New(d.db.ListTasks),
		AddTask:      handler.New(d.client.AddItem),
		TaskDelete:   handler.New(d.client.DeleteItem),
		Sync:         handler.NewPos(d.client.Sync, "isForce"),
		ListProjects: handler.New(d.db.ListProjects),
	})

	if err := server.Loop(ctx, server.NetAccepter(lst, channel.Line), svc, &server.LoopOptions{
		ServerOptions: &jrpc2.ServerOptions{
			Logger: jrpc2.StdLogger(d.log),
		},
	}); err != nil {
		return err
	}

	return nil
}
