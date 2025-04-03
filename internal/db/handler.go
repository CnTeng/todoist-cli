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
		return db.handleSyncResponse(r)
	case *sync.CompletedGetResponse:
		return db.handleCompletedGetResponse(r)
	}

	return nil
}

func (db *DB) handleSyncResponse(resp *sync.SyncResponse) error {
	if err := db.withTx(func(tx *sql.Tx) error {
		if resp.FullSync {
			if err := db.clean(context.Background(), tx); err != nil {
				return err
			}
		}

		for _, item := range resp.Items {
			if err := db.storeTask(tx, item); err != nil {
				return err
			}
		}

		for _, label := range resp.Labels {
			if err := db.storeLabel(tx, label); err != nil {
				return err
			}
		}

		for _, project := range resp.Projects {
			if err := db.storeProject(tx, project); err != nil {
				return err
			}
		}

		if err := db.storeSyncToken(tx, resp.SyncToken); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	for _, syncErr := range resp.SyncStatus {
		if syncErr != nil {
			return syncErr
		}
	}
	return nil
}

func (db *DB) handleCompletedGetResponse(resp *sync.CompletedGetResponse) error {
	return db.withTx(func(tx *sql.Tx) error {
		for _, item := range resp.Items {
			if item.ItemObject != nil {
				if err := db.storeTask(tx, item.ItemObject); err != nil {
					return err
				}
			}
		}

		for _, project := range resp.Projects {
			if err := db.storeProject(tx, project); err != nil {
				return err
			}
		}

		return nil
	})
}
