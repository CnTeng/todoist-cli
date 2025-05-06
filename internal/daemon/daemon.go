package daemon

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/CnTeng/todoist-api-go/todoist"
	"github.com/CnTeng/todoist-api-go/ws"
	"github.com/CnTeng/todoist-cli/internal/db"
	"github.com/creachadair/jrpc2"
	"github.com/creachadair/jrpc2/channel"
	"github.com/creachadair/jrpc2/handler"
	"github.com/creachadair/jrpc2/server"
)

type Config struct {
	Address      string `toml:"address"`
	ApiToken     string `toml:"api_token"`
	ApiTokenFile string `toml:"api_token_file"`
}

var DefaultConfig = &Config{
	Address: "@todo.sock",
}

type Daemon struct {
	address string
	client  *todoist.Client
	db      *db.DB
	log     *log.Logger
}

func NewDaemon(db *db.DB, config *Config) *Daemon {
	return &Daemon{
		address: config.Address,
		client:  todoist.NewClient(http.DefaultClient, config.ApiToken, db),
		db:      db,
		log:     log.New(log.Writer(), "daemon: ", log.Flags()),
	}
}

func (d *Daemon) loadWsToken(ctx context.Context) (string, error) {
	user, err := d.db.GetUser(ctx)
	if err == nil {
		return user.WebsocketURL, nil
	}

	user, err = d.client.GetUser(ctx)
	if err != nil {
		return "", err
	}

	return user.WebsocketURL, nil
}

func (d *Daemon) Serve(ctx context.Context) error {
	url, err := d.loadWsToken(ctx)
	if err != nil {
		return err
	}

	ws := ws.NewClient(url, d)
	ws.Listen(ctx)

	lst, err := net.Listen(jrpc2.Network(d.address))
	if err != nil {
		return err
	}
	defer lst.Close()

	svc := server.Static(handler.Map{
		Sync: handler.New(d.sync),
		// TODO: completed tasks
		TaskGet:       handler.New(d.db.GetTask),
		TaskList:      handler.NewPos(d.db.ListTasks, "completed"),
		TaskAdd:       handler.New(d.client.AddTask),
		TaskQuickAdd:  handler.New(d.client.AddTaskQuick),
		TaskModify:    handler.New(d.client.UpdateTask),
		TaskRemove:    handler.New(d.client.DeleteTasks),
		TaskClose:     handler.New(d.client.CloseTask),
		TaskMove:      handler.New(d.client.MoveTask),
		TaskReopen:    handler.New(d.client.UncompleteTask),
		ProjectGet:    handler.New(d.db.GetProject),
		ProjectList:   handler.New(d.db.ListProjects),
		ProjectAdd:    handler.New(d.client.AddProject),
		ProjectModify: handler.New(d.client.UpdateProject),
		ProjectRemove: handler.New(d.client.DeleteProjects),
		LabelGet:      handler.New(d.db.GetLabel),
		LabelList:     handler.New(d.db.ListLabels),
	})

	return server.Loop(ctx, server.NetAccepter(lst, channel.Line), svc, &server.LoopOptions{
		ServerOptions: &jrpc2.ServerOptions{
			Logger: jrpc2.StdLogger(d.log),
		},
	})
}
