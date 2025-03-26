package db

import (
	"encoding/json"

	"github.com/CnTeng/todoist-api-go/sync/v9"
)

const (
	storeLabelQuery     = `INSERT INTO labels (id, data) VALUES (?, ?) ON CONFLICT (id) DO UPDATE SET data = excluded.data`
	getLabelQueryByID   = `SELECT data FROM labels WHERE id = ?`
	getLabelQueryByName = `SELECT data FROM labels WHERE data ->> 'name' = ?`
)

func (db *DB) StoreLabel(label *sync.Label) error {
	data, err := json.Marshal(label)
	if err != nil {
		return err
	}

	if _, err := db.Exec(storeLabelQuery, label.ID, data); err != nil {
		return err
	}

	return nil
}

func (db *DB) GetLabelByID(id string) (*sync.Label, error) {
	var data []byte
	if err := db.QueryRow(getLabelQueryByID, id).Scan(&data); err != nil {
		return nil, err
	}

	label := &sync.Label{}
	if err := json.Unmarshal(data, label); err != nil {
		return nil, err
	}

	return label, nil
}

func (db *DB) GetLabelByName(name string) (*sync.Label, error) {
	var data []byte
	if err := db.QueryRow(getLabelQueryByName, name).Scan(&data); err != nil {
		return nil, err
	}

	label := &sync.Label{}
	if err := json.Unmarshal(data, label); err != nil {
		return nil, err
	}

	return label, nil
}
