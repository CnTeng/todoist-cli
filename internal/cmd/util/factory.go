package util

import (
	"net"

	"github.com/CnTeng/todoist-cli/internal/cli"
	"github.com/CnTeng/todoist-cli/internal/daemon"
	"github.com/creachadair/jrpc2"
	"github.com/creachadair/jrpc2/channel"
)

type Factory struct {
	RpcClient    *jrpc2.Client
	DeamonConfig *daemon.Config `toml:"daemon"`
	Lang         string         `toml:"lang"`
	Cli          *cli.Cli

	conn net.Conn
}

func NewFactory() *Factory {
	return &Factory{
		RpcClient:    nil,
		DeamonConfig: &daemon.Config{},
		Lang:         "en",
		Cli:          cli.NewCLI(cli.Nerd),
	}
}

func (f *Factory) Dial() error {
	var err error
	f.conn, err = net.Dial("unix", "@todo.sock")
	if err != nil {
		return err
	}

	f.RpcClient = jrpc2.NewClient(channel.Line(f.conn, f.conn), nil)
	return nil
}

func (f *Factory) Close() error {
	if f.conn != nil {
		return f.conn.Close()
	}
	return nil
}
