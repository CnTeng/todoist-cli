package daemon

import (
	"context"
	"encoding/json"
	"log"

	"github.com/creachadair/jrpc2"
)

type rpcLogger struct {
	logger *log.Logger
}

func (rl *rpcLogger) LogRequest(ctx context.Context, req *jrpc2.Request) {
	reqType := "request"
	if req.IsNotification() {
		reqType = "notification"
	}

	id := req.ID()
	if id == "" {
		id = "-"
	}

	var params json.RawMessage
	_ = req.UnmarshalParams(&params)
	if len(params) == 0 {
		params = json.RawMessage(`null`)
	}

	rl.logger.Printf("%s for %s (ID %s): %s", reqType, req.Method(), id, params)
}

func (rl *rpcLogger) LogResponse(ctx context.Context, rsp *jrpc2.Response) {
	id := rsp.ID()
	if id == "" {
		id = "-"
	}

	req := jrpc2.InboundRequest(ctx)

	if rsp.Error() == nil {
		rl.logger.Printf("response to %s (ID %s): success", req.Method(), id)
	} else {
		rl.logger.Printf("response to %s (ID %s): failed, %s", req.Method(), id, rsp.Error())
	}
}
