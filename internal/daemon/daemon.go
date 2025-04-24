package daemon

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/CnTeng/todoist-api-go/todoist"
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
	client  *todoist.Client
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

	d.client = todoist.NewClient(http.DefaultClient, apiToken, d.db)
	d.ws = ws.NewClient(wsToken, d)

	return nil
}

func (d *Daemon) Serve(ctx context.Context) error {
	d.ws.Listen(ctx)

	lst, err := net.Listen(jrpc2.Network(d.address))
	if err != nil {
		return err
	}
	defer lst.Close()

	svc := server.Static(handler.Map{
		Sync: handler.New(d.sync),
		// TODO: completed tasks
		CompletedGet:  handler.New(d.client.GetTasksCompletedByDueDate),
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
