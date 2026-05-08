package monitor

import (
	"github.com/go-ping/ping"
)

func ICMPCheck(address string) (string, int64) {
	pinger, err := ping.NewPinger(address)
	if err != nil {
		return address, -1
	}

	pinger.Count = 1
	pinger.Timeout = 1000000000 // 1 second in nanoseconds

	err = pinger.Run()
	if err != nil {
		return address, -1
	}

	stats := pinger.Statistics()
	if stats.PacketsRecv == 0 {
		return address, -1
	}

	return address, stats.AvgRtt.Milliseconds()
}
