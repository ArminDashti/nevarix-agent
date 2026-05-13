package monitor

import (
	"encoding/json"
	"os"
	"strconv"
)

// ConfigFile is the path to the monitor loop JSON configuration.
const ConfigFile = "/home/.nevarix-server/monitor.json"

const runtimeConfigPath = "/home/.nevarix-server/runtime.json"

// LoadConfig reads monitor settings; missing file yields defaults.
func LoadConfig(path string) (*Config, error) {
	cfg := &Config{IntervalSeconds: 60}
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			if v := os.Getenv("AGENT_MONITOR_INTERVAL_SECONDS"); v != "" {
				if n, e := strconv.Atoi(v); e == nil && n > 0 {
					cfg.IntervalSeconds = n
				}
			}
			return cfg, nil
		}
		return nil, err
	}
	if err := json.Unmarshal(b, cfg); err != nil {
		return nil, err
	}
	if cfg.IntervalSeconds <= 0 {
		cfg.IntervalSeconds = 60
	}
	return cfg, nil
}

// GetRuntimeConfig loads hub URL from file or environment.
func GetRuntimeConfig() (*RuntimeConfig, error) {
	rc := &RuntimeConfig{}
	if env := os.Getenv("NEVARIX_HUB_BASE_URL"); env != "" {
		rc.API.HubBaseURL = env
		return rc, nil
	}
	b, err := os.ReadFile(runtimeConfigPath)
	if err != nil {
		return rc, err
	}
	if err := json.Unmarshal(b, rc); err != nil {
		return nil, err
	}
	return rc, nil
}
