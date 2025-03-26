CREATE TABLE IF NOT EXISTS items (id text PRIMARY KEY, data jsonb NOT NULL);

CREATE TABLE IF NOT EXISTS projects (id text PRIMARY KEY, data jsonb NOT NULL);

CREATE TABLE IF NOT EXISTS labels (id text PRIMARY KEY, data jsonb NOT NULL);

CREATE TABLE IF NOT EXISTS sync_token (
  id integer PRIMARY KEY CHECK (id = 1),
  token text NOT NULL,
  last_sync DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- INSERT INTO
--   sync_token (id, token)
-- VALUES
--   (1, "*")
-- ON CONFLICT (id) DO UPDATE
-- SET
--   token = excluded.token;
