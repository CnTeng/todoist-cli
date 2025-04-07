package util

import (
	"net"

	"github.com/CnTeng/todoist-cli/internal/daemon"
	"github.com/CnTeng/todoist-cli/internal/db"
	"github.com/CnTeng/todoist-cli/internal/view"
	"github.com/creachadair/jrpc2"
	"github.com/creachadair/jrpc2/channel"
)

type Factory struct {
	DeamonConfig *daemon.Config   `toml:"daemon"`
	ViewConfig   *view.ViewConfig `toml:"view"`
	DBConfig     *db.Config       `toml:"database"`
	IconConfig   *view.IconConfig `toml:"icon"`

	Lang string `toml:"lang"`

	conn net.Conn
	*jrpc2.Client
}

func NewFactory() *Factory {
	return &Factory{
		DeamonConfig: daemon.DefaultConfig,
		ViewConfig:   view.DefaultViewConfig,
		IconConfig:   view.DefaultIconConfig,

		Lang: "en",
	}
}

func (f *Factory) Dial() error {
	var err error
	f.conn, err = net.Dial("unix", "@todo.sock")
	if err != nil {
		return err
	}

	f.Client = jrpc2.NewClient(channel.Line(f.conn, f.conn), nil)
	return nil
}

func (f *Factory) Close() error {
	if f.conn != nil {
		return f.conn.Close()
	}
	return nil
}
