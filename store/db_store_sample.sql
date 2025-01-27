CREATE TABLE IF NOT EXISTS passwordless_tokens (
	id TEXT PRIMARY KEY,
	recipient TEXT NOT NULL,
	code_hash BLOB NOT NULL,
	expires_at DATETIME NOT NULL,
	created_at DATETIME NOT NULL,
	attempts INTEGER NOT NULL DEFAULT 0
  );