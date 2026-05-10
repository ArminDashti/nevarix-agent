package database

import (
	"database/sql"
	"time"
)

// InsertHttpLatency inserts or updates HTTP latency data.
func InsertHttpLatency(db *sql.DB, endpointAddress string, latencyMS int64, checkedAt time.Time) error {
	_, err := db.Exec(
		`INSERT INTO http_latencies (endpoint_address, latency_ms, checked_at) 
		 VALUES (?, ?, ?) 
		 ON CONFLICT(endpoint_address, checked_at) DO UPDATE SET 
		 latency_ms = excluded.latency_ms`,
		endpointAddress,
		latencyMS,
		checkedAt.Unix(),
	)
	return err
}

