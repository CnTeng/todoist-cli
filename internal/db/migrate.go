package db

import "database/sql"

const createTableQuery = `
	CREATE TABLE IF NOT EXISTS tasks (id text PRIMARY KEY, data jsonb NOT NULL);

	CREATE TABLE IF NOT EXISTS projects (id text PRIMARY KEY, data jsonb NOT NULL);

	CREATE TABLE IF NOT EXISTS labels (id text PRIMARY KEY, data jsonb NOT NULL);

	CREATE TABLE IF NOT EXISTS users (id text PRIMARY KEY, data jsonb NOT NULL);

	CREATE TABLE IF NOT EXISTS sync_token (
		id integer PRIMARY KEY CHECK (id = 1),
		token text NOT NULL,
		last_sync DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	INSERT INTO
		sync_token (id, token)
	VALUES
		(1, "*")
	ON CONFLICT (id) DO NOTHING;`

const createViewQuery = `
	CREATE VIEW IF NOT EXISTS labels_view AS
	WITH
		label_count AS (
			SELECT
				je.value AS id,
				count(tasks.id) AS count
			FROM
				tasks,
				json_each(tasks.data -> 'labels') AS je
			GROUP BY
				je.value
		)
	SELECT
		l.id,
		json_patch(
			l.data,
			json_object('count', coalesce(lc.count, 0))
		) AS data
	FROM
		labels l
		LEFT JOIN label_count lc ON l.id = lc.id
	UNION ALL
	SELECT
		lc.id AS id,
		json_object(
			'name',
			lc.id,
			'is_shared',
			json('true'),
			'count',
			lc.count
		) AS data
	FROM
		label_count lc
	WHERE
		NOT EXISTS (
			SELECT
				1
			FROM
				labels l
			WHERE
				l.id = lc.id
		);`

var migrations = []string{
	createTableQuery,
	createViewQuery,
}

func (db *DB) Migrate() error {
	return db.withTx(func(tx *sql.Tx) error {
		for _, migration := range migrations {
			if _, err := tx.Exec(migration); err != nil {
				return err
			}
		}
		return nil
	})
}
