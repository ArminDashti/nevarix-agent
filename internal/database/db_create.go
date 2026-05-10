package database

import (
	"database/sql"
)

// EnsureSchema creates the necessary tables for SQLite if they do not already exist.
func EnsureSchema(db *sql.DB) error {
	// Table for HTTP latency results
	// Columns: endpoint_address, latency_ms, checked_at
	queryHttp := `
	CREATE TABLE IF NOT EXISTS http_latencies (
		endpoint_address TEXT,
		latency_ms INTEGER,
		checked_at INTEGER,
		UNIQUE(endpoint_address, checked_at)
	);`

	if _, err := db.Exec(queryHttp); err != nil {
		return err
	}

	return nil
}
