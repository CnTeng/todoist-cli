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

	user, err = todoist.NewUserService(d.client).GetUser(ctx)
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

	taskSvc := todoist.NewTaskService(d.client)
	projectSvc := todoist.NewProjectService(d.client)
	sectionSvc := todoist.NewSectionService(d.client)
	labelSvc := todoist.NewLabelService(d.client)
	filterSvc := todoist.NewFilterService(d.client)

	svc := server.Static(handler.Map{
		Sync: handler.New(d.sync),

		// Task services
		TaskList:     handler.New(d.db.ListTasks),
		TaskAdd:      handler.New(taskSvc.AddTask),
		TaskQuickAdd: handler.New(taskSvc.QuickAddTask),
		TaskModify:   handler.New(taskSvc.UpdateTask),
		TaskMove:     handler.New(taskSvc.MoveTask),
		TaskReorder:  handler.New(taskSvc.ReorderTasks),
		TaskClose:    handler.New(taskSvc.CloseTasks),
		TaskReopen:   handler.New(taskSvc.UncompleteTasks),
		TaskRemove:   handler.New(taskSvc.DeleteTasks),

		// Project services
		ProjectList:      handler.New(d.db.ListProjects),
		ProjectAdd:       handler.New(projectSvc.AddProject),
		ProjectModify:    handler.New(projectSvc.UpdateProject),
		ProjectReorder:   handler.New(projectSvc.ReorderProjects),
		ProjectArchive:   handler.New(projectSvc.ArchiveProjects),
		ProjectUnarchive: handler.New(projectSvc.UnarchiveProjects),
		ProjectRemove:    handler.New(projectSvc.DeleteProjects),

		// Section services
		SectionList:      handler.New(d.db.ListSections),
		SectionAdd:       handler.New(sectionSvc.AddSection),
		SectionModify:    handler.New(sectionSvc.UpdateSection),
		SectionMove:      handler.New(sectionSvc.MoveSection),
		SectionReorder:   handler.New(sectionSvc.ReorderSections),
		SectionArchive:   handler.New(sectionSvc.ArchiveSections),
		SectionUnarchive: handler.New(sectionSvc.UnarchiveSections),
		SectionRemove:    handler.New(sectionSvc.DeleteSections),

		// Label services
		LabelList:    handler.New(d.db.ListLabels),
		LabelAdd:     handler.New(labelSvc.AddLabel),
		LabelModify:  handler.New(d.updateLabel),
		LabelReorder: handler.New(labelSvc.ReorderLabels),
		LabelRemove:  handler.New(d.deleteLabels),

		// Filter services
		FilterList:    handler.New(d.db.ListFilters),
		FilterAdd:     handler.New(filterSvc.AddFilter),
		FilterModify:  handler.New(filterSvc.UpdateFilter),
		FilterReorder: handler.New(filterSvc.ReorderFilters),
		FilterRemove:  handler.New(filterSvc.DeleteFilters),
	})

	return server.Loop(ctx, server.NetAccepter(lst, channel.Line), svc, &server.LoopOptions{
		ServerOptions: &jrpc2.ServerOptions{
			Logger: jrpc2.StdLogger(d.log),
		},
	})
}
