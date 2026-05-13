package monitor

// MonitorState tracks a detached agent child process.
type MonitorState struct {
	PID           int   `json:"pid"`
	StartedAtUnix int64 `json:"started_at_unix"`
}

// Config controls the monitoring loop.
type Config struct {
	IntervalSeconds int `json:"interval_seconds"`
}

// RuntimeConfig holds hub and API wiring loaded from disk or environment.
type RuntimeConfig struct {
	API struct {
		HubBaseURL string `json:"hub_base_url"`
	} `json:"api"`
}
