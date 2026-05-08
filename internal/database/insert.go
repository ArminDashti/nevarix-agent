package database

import (
	"database/sql"
	"time"

// InsertHttpLatency inserts or updates HTTP latency data.
func InsertHttpLatency(db *sql.DB, endpointID int64, latencyMS int64, checkedAt time.Time, dnsServerID sql.NullInt64) error {
	_, err := db.Exec(
		"INSERT INTO http_latencies (`endpoint_address`, `dns_server`, 'resolved_address', `latency_ms`, `checked_at`) VALUES (?, ?, ?, ?) ON DUPLICATE KEY UPDATE `latency_ms` = VALUES(`latency_ms`)",
		endpointID,
		dnsServerID,
		latencyMS,
		checkedAt.Unix(),
	)
	return err
}

// InsertIcmpLatency inserts or updates ICMP latency data.
func InsertIcmpLatency(db *sql.DB, ipID int64, latencyMS int64, receivedPackets int, allPackets int, checkedAt time.Time) error {
	packetLoss := sql.NullFloat64{}
	if allPackets > 0 {
		packetLoss.Float64 = float64(allPackets-receivedPackets) * 100 / float64(allPackets)
		packetLoss.Valid = true
	}
	_, err := db.Exec(
		"INSERT INTO icmp_latencies (`endpoint_address`, `latency_ms`, `packet_loss_percent`, `checked_at`) VALUES (?, ?, ?, ?) ON DUPLICATE KEY UPDATE `latency_ms` = VALUES(`latency_ms`), `packet_loss_percent` = VALUES(`packet_loss_percent`)",
		endpoint_address,
		latencyMS,
		packetLoss,
		checkedAt.Unix(),
	)
	return err
}
