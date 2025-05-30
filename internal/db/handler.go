package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/CnTeng/todoist-api-go/rest"
	"github.com/CnTeng/todoist-api-go/sync"
)

func (db *DB) ResourceTypes(ctx context.Context) (*sync.ResourceTypes, error) {
	return &sync.ResourceTypes{sync.All}, nil
}

func (db *DB) SyncToken(ctx context.Context) (*string, error) {
	return db.getSyncToken(ctx)
}

func (db *DB) HandleResponse(ctx context.Context, resp any) error {
	switch r := resp.(type) {
	case *sync.SyncResponse:
		return db.handleSyncResponse(ctx, r)
	case *rest.TaskGetCompletedResponse:
		return db.handleTaskGetCompletedResponse(ctx, r)
	case *rest.ProjectGetArchivedResponse:
		return db.handleProjectGetArchivedResponse(ctx, r)
	}

	return nil
}

func (db *DB) handleSyncResponse(ctx context.Context, resp *sync.SyncResponse) error {
	if err := db.withTx(func(tx *sql.Tx) error {
		if resp.FullSync {
			if err := db.clean(ctx, tx); err != nil {
				return err
			}
		}

		for _, task := range resp.Tasks {
			if err := db.storeTask(ctx, tx, task); err != nil {
				return fmt.Errorf("failed to store task %s: %w", task.ID, err)
			}
		}

		for _, project := range resp.Projects {
			if err := db.storeProject(ctx, tx, project); err != nil {
				return fmt.Errorf("failed to store project %s: %w", project.ID, err)
			}
		}

		for _, section := range resp.Sections {
			if err := db.storeSection(ctx, tx, section); err != nil {
				return fmt.Errorf("failed to store section %s: %w", section.ID, err)
			}
		}

		for _, label := range resp.Labels {
			if err := db.storeLabel(ctx, tx, label); err != nil {
				return fmt.Errorf("failed to store label %s: %w", label.ID, err)
			}
		}

		for _, filter := range resp.Filters {
			if err := db.storeFilter(ctx, tx, filter); err != nil {
				return fmt.Errorf("failed to store filter %s: %w", filter.ID, err)
			}
		}

		if resp.User != nil {
			if err := db.storeUser(ctx, tx, resp.User); err != nil {
				return fmt.Errorf("failed to store user %s: %w", resp.User.ID, err)
			}
		}

		if err := db.storeSyncToken(ctx, tx, resp.SyncToken); err != nil {
			return fmt.Errorf("failed to store sync token: %w", err)
		}

		return nil
	}); err != nil {
		return err
	}

	for _, syncErr := range resp.SyncStatus {
		if syncErr != nil {
			return fmt.Errorf("sync error: %w", syncErr)
		}
	}
	return nil
}

func (db *DB) handleTaskGetCompletedResponse(ctx context.Context, resp *rest.TaskGetCompletedResponse) error {
	return db.withTx(func(tx *sql.Tx) error {
		for _, task := range resp.Tasks {
			if err := db.storeTask(ctx, tx, task); err != nil {
				return fmt.Errorf("failed to store completed task %s: %w", task.ID, err)
			}
		}

		return nil
	})
}

func (db *DB) handleProjectGetArchivedResponse(ctx context.Context, resp *rest.ProjectGetArchivedResponse) error {
	return db.withTx(func(tx *sql.Tx) error {
		for _, project := range resp.Projects {
			if err := db.storeProject(ctx, tx, project); err != nil {
				return fmt.Errorf("failed to store archived project %s: %w", project.ID, err)
			}
		}

		return nil
	})
}
