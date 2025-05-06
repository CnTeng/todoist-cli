package util

import (
	"net"
	"path/filepath"

	"github.com/CnTeng/todoist-cli/internal/daemon"
	"github.com/CnTeng/todoist-cli/internal/view"
	"github.com/adrg/xdg"
	"github.com/creachadair/jrpc2"
	"github.com/creachadair/jrpc2/channel"
)

const (
	appName    = "todoist"
	configFile = "config.toml"
	dataFile   = "todoist.db"
)

type Factory struct {
	DeamonConfig *daemon.Config   `toml:"daemon"`
	IconConfig   *view.IconConfig `toml:"icon"`

	ConfigPath string `toml:"-"`
	DataPath   string `toml:"-"`

	conn net.Conn
	*jrpc2.Client
}

func NewFactory() (*Factory, error) {
	configPath, err := xdg.ConfigFile(filepath.Join(appName, configFile))
	if err != nil {
		return nil, err
	}
	dataPath, err := xdg.DataFile(filepath.Join(appName, dataFile))
	if err != nil {
		return nil, err
	}

	return &Factory{
		DeamonConfig: daemon.DefaultConfig,
		IconConfig:   view.DefaultIconConfig,

		ConfigPath: configPath,
		DataPath:   dataPath,
	}, nil
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
