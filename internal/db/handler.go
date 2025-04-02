package db

import (
	"context"
	"database/sql"

	"github.com/CnTeng/todoist-api-go/sync"
)

func (db *DB) ResourceTypes() (*sync.ResourceTypes, error) {
	return &sync.ResourceTypes{sync.All}, nil
}

func (db *DB) SyncToken() (*string, error) {
	return db.getSyncToken()
}

func (db *DB) HandleResponse(resp any) error {
	switch r := resp.(type) {
	case *sync.SyncResponse:
		return db.withTx(func(tx *sql.Tx) error {
			if r.FullSync {
				if err := db.clean(context.Background(), tx); err != nil {
					return err
				}
			}

			for _, item := range r.Items {
				if err := db.storeTask(tx, item); err != nil {
					return err
				}
			}

			for _, label := range r.Labels {
				if err := db.storeLabel(tx, label); err != nil {
					return err
				}
			}

			for _, project := range r.Projects {
				if err := db.storeProject(tx, project); err != nil {
					return err
				}
			}

			if err := db.storeSyncToken(tx, r.SyncToken); err != nil {
				return err
			}

			return nil
		})
	}

	return nil
}
