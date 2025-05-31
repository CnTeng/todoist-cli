package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
)

func getItem[T any](ctx context.Context, tx *sql.Tx, query string, args ...any) (*T, error) {
	var data []byte
	if err := tx.QueryRowContext(ctx, query, args...).Scan(&data); err != nil {
		return nil, fmt.Errorf("failed to get item: %w", err)
	}

	var item T
	if err := json.Unmarshal(data, &item); err != nil {
		return nil, fmt.Errorf("failed to unmarshal item: %w", err)
	}

	return &item, nil
}

func listItems[T any](ctx context.Context, tx *sql.Tx, query string, args ...any) ([]*T, error) {
	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []*T{}
	for rows.Next() {
		var data []byte
		if err := rows.Scan(&data); err != nil {
			return nil, fmt.Errorf("failed to list items: %w", err)
		}

		var item T
		if err := json.Unmarshal(data, &item); err != nil {
			return nil, fmt.Errorf("failed to unmarshal item: %w", err)
		}

		items = append(items, &item)
	}

	return items, nil
}
