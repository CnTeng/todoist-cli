package util

import (
	"github.com/CnTeng/todoist-cli/internal/model"
	"github.com/creachadair/jrpc2"
)

type Factory struct {
	RpcClient *jrpc2.Client
	Config    *model.Config
}
