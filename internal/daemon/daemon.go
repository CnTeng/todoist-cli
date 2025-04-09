package daemon

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/CnTeng/todoist-api-go/sync"
	"github.com/CnTeng/todoist-api-go/ws"
	"github.com/CnTeng/todoist-cli/internal/db"
	"github.com/creachadair/jrpc2"
	"github.com/creachadair/jrpc2/channel"
	"github.com/creachadair/jrpc2/handler"
	"github.com/creachadair/jrpc2/server"
)

const (
	apiTokenEnv     = "API_TOKEN"
	apiTokenFileEnv = "API_TOKEN_FILE"
	WsTokenEnv      = "WS_TOKEN"
	WsTokenFileEnv  = "WS_TOKEN_FILE"
)

type Config struct {
	Address string `toml:"address"`
}

var DefaultConfig = &Config{
	Address: "@todo.sock",
}

type Daemon struct {
	address string
	client  *sync.Client
	ws      *ws.Client
	db      *db.DB
	log     *log.Logger
}

func NewDaemon(db *db.DB, config *Config) *Daemon {
	d := &Daemon{
		address: config.Address,
		db:      db,
		log:     log.New(log.Writer(), "daemon: ", log.Flags()),
	}

	return d
}

func (d *Daemon) LoadTokens() error {
	apiToken := os.Getenv(apiTokenEnv)
	apiTokenFile := os.Getenv(apiTokenFileEnv)
	if apiToken == "" && apiTokenFile == "" {
		return fmt.Errorf("%s or %s is required", apiTokenEnv, apiTokenFileEnv)
	}
	if apiToken == "" && apiTokenFile != "" {
		token, err := os.ReadFile(apiTokenFile)
		if err != nil {
			return err
		}
		apiToken = string(token)
	}

	wsToken := os.Getenv(WsTokenEnv)
	wsTokenFile := os.Getenv(WsTokenFileEnv)
	if wsToken == "" && wsTokenFile == "" {
		return fmt.Errorf("%s or %s is required,", WsTokenEnv, WsTokenFileEnv)
	}
	if wsToken == "" && wsTokenFile != "" {
		token, err := os.ReadFile(wsTokenFile)
		if err != nil {
			return err
		}
		wsToken = string(token)
	}

	d.client = sync.NewClient(http.DefaultClient, apiToken, d.db)
	d.ws = ws.NewClient(wsToken, d)

	return nil
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
		TaskList:      handler.NewPos(d.db.ListTasks, "completed"),
		TaskAdd:       handler.New(d.client.AddItem),
		TaskQuickAdd:  handler.New(d.client.QuickAddItem),
		TaskModify:    handler.New(d.client.UpdateItem),
		TaskRemove:    handler.New(d.client.DeleteItems),
		TaskClose:     handler.New(d.client.CloseItem),
		TaskMove:      handler.New(d.client.MoveItem),
		TaskReopen:    handler.New(d.client.UncompleteItem),
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
