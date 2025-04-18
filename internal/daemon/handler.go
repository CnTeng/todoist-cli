package daemon

import (
	"context"

	"github.com/CnTeng/todoist-api-go/sync"
	"github.com/CnTeng/todoist-api-go/ws"
)

func (d *Daemon) HandleMessage(ctx context.Context, msg ws.Message) error {
	switch msg {
	case ws.Connected:
		fallthrough
	case ws.SyncNeeded:
		d.log.Println("sync needed")
		if _, err := d.client.Sync(ctx, false); err != nil {
			return err
		}
	}
	return nil
}

type SyncArgs struct {
	IsForce bool `json:"isForce"`
	All     bool `json:"all"`
}

func (d *Daemon) sync(ctx context.Context, args *SyncArgs) error {
	if _, err := d.client.Sync(ctx, args.IsForce); err != nil {
		return err
	}

	if !args.All {
		return nil
	}

	ps, err := d.db.ListProjects(ctx)
	if err != nil {
		return err
	}

	for _, p := range ps {
		AnnotateItems := true
		params := &sync.CompletedGetParams{
			ProjectID:     &p.ID,
			AnnotateItems: &AnnotateItems,
		}

		if _, err := d.client.GetCompletedInfo(ctx, params); err != nil {
			return err
		}
	}

	return nil
}
