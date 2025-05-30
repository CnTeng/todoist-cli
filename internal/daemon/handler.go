package daemon

import (
	"context"

	"github.com/CnTeng/todoist-api-go/ws"
)

func (d *Daemon) HandleMessage(ctx context.Context, msg ws.Message) error {
	switch msg {
	case ws.Connected:
		fallthrough
	case ws.SyncNeeded:
		if _, err := d.client.SyncWithAutoToken(ctx, false); err != nil {
			return err
		}
	}
	return nil
}
