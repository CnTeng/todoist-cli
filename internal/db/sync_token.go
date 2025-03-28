package db

const (
	storeSyncTokenQuery = `INSERT INTO sync_token (id, token) VALUES (1, ?) ON CONFLICT (id) DO UPDATE SET token = excluded.token`
	getSyncTokenQuery   = `SELECT token FROM sync_token WHERE id = 1`
)

func (db *DB) storeSyncToken(token string) error {
	if _, err := db.Exec(storeSyncTokenQuery, token); err != nil {
		return err
	}

	return nil
}

func (db *DB) getSyncToken() (*string, error) {
	var token string

	if err := db.QueryRow(getSyncTokenQuery).Scan(&token); err != nil {
		return nil, err
	}

	return &token, nil
}
