package db

import "github.com/CnTeng/todoist-api-go/sync/v9"

func (db *DB) ResourceTypes() (*sync.ResourceTypes, error) {
	return &sync.ResourceTypes{sync.All}, nil
}

func (db *DB) SyncToken() (*string, error) {
	return db.getSyncToken()
}

func (db *DB) HandleResponse(resp any) error {
	switch r := resp.(type) {
	case *sync.SyncResponse:
		for _, item := range r.Items {
			if err := db.StoreTask(item); err != nil {
				return err
			}
		}

		for _, label := range r.Labels {
			if err := db.StoreLabel(label); err != nil {
				return err
			}
		}

		for _, project := range r.Projects {
			if err := db.StoreProject(project); err != nil {
				return err
			}
		}

		if err := db.storeSyncToken(r.SyncToken); err != nil {
			return err
		}
	}

	return nil
}
