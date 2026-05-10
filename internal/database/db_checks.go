package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// CheckAddressesAndStore checks configured endpoints and stores each result.
func CheckAddressesAndStore(connStr string, cfg Config, checkedAt time.Time) error {
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		return err
	}
	defer db.Close()
	if err := EnsureSchema(db); err != nil {
		return err
	}
	type1Endpoints, err := loadEndpointsByType(db, 1)
	if err != nil {
		return err
	}

	roundedCheckedAt := checkedAt.Truncate(time.Minute)
	if len(type1Endpoints) > 0 {
		checkPort := 443
		if len(cfg.Ports) > 0 && cfg.Ports[0] > 0 {
			checkPort = cfg.Ports[0]
		}
		for _, endpoint := range type1Endpoints {
			latencyMS, _, _ := probeEndpoint(db, endpoint.ID, endpoint.Host, checkPort, time.Duration(endpointTimeoutMSDefault())*time.Millisecond)
			if err := upsertLatency(db, endpoint.ID, latencyMS, roundedCheckedAt); err != nil {
				writeErrorLog(err)
				return err
			}
		}
	}
	return nil
}