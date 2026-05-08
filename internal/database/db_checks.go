package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// CheckAddressesAndStore checks configured hosts and stores each result.
func CheckAddressesAndStore(connStr string, cfg Config, checkedAt time.Time) error {
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		return err
	}
	defer db.Close()
	if err := EnsureSchema(db); err != nil {
		return err
	}
	type1Hosts, err := loadHostsByType(db, 1)
	if err != nil {
		return err
	}
	dnsServers, err := loadActiveDNSServers(db)
	if err != nil {
		return err
	}

	roundedCheckedAt := checkedAt.Truncate(time.Minute)
	if len(type1Hosts) > 0 {
		checkPort := 443
		if len(cfg.Ports) > 0 && cfg.Ports[0] > 0 {
			checkPort = cfg.Ports[0]
		}
		for _, endpoint := range type1Hosts {
			latencyMS, _, dnsServerID := probeHostViaDNSServers(db, endpoint.ID, endpoint.Host, checkPort, dnsServers, time.Duration(hostTimeoutMSDefault())*time.Millisecond, time.Duration(dnsTimeoutMSDefault())*time.Millisecond)
			if err := upsertLatency(db, endpoint.ID, latencyMS, roundedCheckedAt, dnsServerID); err != nil {
				writeErrorLog(err)
				return err
			}
		}
	}

	for _, dnsServer := range dnsServers {
		startAt := time.Now().Local()
		latencyMS, receivedPackets, allPackets, pingErr := probeICMPLatencyAverage(dnsServer.Address, 4, time.Duration(hostTimeoutMSDefault())*time.Millisecond)
		if pingErr != nil {
			latencyMS = -1
			writeProbeLogLine(startAt, "PING", dnsServer.Address, "", "FAIL", time.Now().Local(), "TIMEOUT", "TIMEOUT")
			writeErrorLog(pingErr)
		} else {
			writeProbeLogLine(startAt, "PING", dnsServer.Address, "", "SUCCESS", time.Now().Local(), fmt.Sprintf("%d", latencyMS), "")
		}
		if err := upsertDNSLatency(db, dnsServer.IPID, latencyMS, receivedPackets, allPackets, roundedCheckedAt); err != nil {
			writeErrorLog(err)
			return err
		}
	}
	return nil
}