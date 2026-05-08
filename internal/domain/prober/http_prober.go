package monitor

import (
	"net/http"
	"strings"
	"time"
)

func AddressAvailibility(address string) int64 {
	requestURL := strings.TrimSpace(address)

	client := http.Client{
		Timeout: 10 * time.Second,
	}

	start := time.Now()
	resp, err := client.Get(requestURL)
	if err != nil {
		return -1
	}
	defer resp.Body.Close()
	return time.Since(start).Milliseconds()
}

func DurationUntilNextMinute(now time.Time) time.Duration {
	minuteStart := now.Truncate(time.Minute)
	if now.Equal(minuteStart) {
		return 0
	}
	return minuteStart.Add(time.Minute).Sub(now)
}
