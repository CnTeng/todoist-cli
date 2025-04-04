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

type Config struct {
	Address  string `toml:"address"`
	ApiToken string `toml:"api_token"`
	WsToken  string `toml:"ws_token"`
}

type Daemon struct {
	address string
	client  *sync.Client
	ws      *ws.Client
	db      *db.DB
	log     *log.Logger
}

func NewDaemon(config *Config, db *db.DB) *Daemon {
	d := &Daemon{
		address: config.Address,
		client:  sync.NewClient(http.DefaultClient, config.ApiToken, db),
		db:      db,
		log:     log.New(log.Writer(), "daemon: ", log.Flags()),
	}
	d.ws = ws.NewClient(config.WsToken, d)

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
		Sync:          handler.New(d.sync),
		CompletedGet:  handler.New(d.client.GetCompletedInfo),
		TaskGet:       handler.New(d.db.GetTask),
		TaskList:      handler.NewPos(d.db.ListTasks, "all"),
		TaskAdd:       handler.New(d.client.AddItem),
		TaskQuickAdd:  handler.New(d.client.QuickAddItem),
		TaskModify:    handler.New(d.client.UpdateItem),
		TaskRemove:    handler.New(d.client.DeleteItem),
		TaskClose:     handler.New(d.client.CloseItem),
		TaskMove:      handler.New(d.client.MoveItem),
		TaskReopen:    handler.New(d.client.UncompleteItem),
		ProjectGet:    handler.New(d.db.GetProject),
		ProjectList:   handler.New(d.db.ListProjects),
		ProjectAdd:    handler.New(d.client.AddProject),
		ProjectModify: handler.New(d.client.UpdateProject),
		ProjectRemove: handler.New(d.client.DeleteProject),
		LabelGet:      handler.New(d.db.GetLabel),
		LabelList:     handler.New(d.db.ListLabels),
	})

	return server.Loop(ctx, server.NetAccepter(lst, channel.Line), svc, &server.LoopOptions{
		ServerOptions: &jrpc2.ServerOptions{
			Logger: jrpc2.StdLogger(d.log),
		},
	})
}
