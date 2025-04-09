package util

import (
	"net"

	"github.com/CnTeng/todoist-cli/internal/daemon"
	"github.com/CnTeng/todoist-cli/internal/view"
	"github.com/creachadair/jrpc2"
	"github.com/creachadair/jrpc2/channel"
)

type Factory struct {
	DeamonConfig *daemon.Config   `toml:"daemon"`
	IconConfig   *view.IconConfig `toml:"icon"`

	Lang string `toml:"lang"`

	ConfigFilePath string `toml:"-"`
	DataFilePath   string `toml:"-"`

	conn net.Conn
	*jrpc2.Client
}

func NewFactory(configFile, dataFile string) *Factory {
	return &Factory{
		DeamonConfig: daemon.DefaultConfig,
		IconConfig:   view.DefaultIconConfig,
		Lang:         "en",

		ConfigFilePath: configFile,
		DataFilePath:   dataFile,
	}
}

func (f *Factory) Dial() error {
	var err error
	f.conn, err = net.Dial(jrpc2.Network(f.DeamonConfig.Address))
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
