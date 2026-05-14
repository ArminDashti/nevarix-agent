package monitor

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

const defaultAgentDataDir = "/home/.nevarix-server"

// agentDataDir returns the directory used for runtime state and secrets.
// NEVARIX_AGENT_RUNTIME_DIR overrides the default path (for tests or custom installs).
func agentDataDir() string {
	if v := strings.TrimSpace(os.Getenv("NEVARIX_AGENT_RUNTIME_DIR")); v != "" {
		return v
	}
	return defaultAgentDataDir
}

// DefaultMonitorConfigPath is the standard path to optional monitor tuning (JSON).
func DefaultMonitorConfigPath() string {
	return filepath.Join(agentDataDir(), "monitor_config.json")
}

func stateFilePath() string {
	return filepath.Join(agentDataDir(), "agent_state.json")
}

func runtimeFilePath() string {
	return filepath.Join(agentDataDir(), "runtime.json")
}

// apiTokenSecretPath is the filesystem location for the persisted local API bearer token.
func apiTokenSecretPath() string {
	return filepath.Join(agentDataDir(), "api_token.secret")
}

// MonitorState records a detached agent child process.
type MonitorState struct {
	PID           int   `json:"pid"`
	StartedAtUnix int64 `json:"started_at_unix"`
}

// Config holds monitor loop settings.
type Config struct {
	IntervalSeconds int `json:"interval_seconds"`
}

// RuntimeConfig is persisted hub and API settings for the agent.
type RuntimeConfig struct {
	API struct {
		HubBaseURL string `json:"hub_base_url"`
	} `json:"api"`
}

// EnsureRuntimeIntegrity ensures the agent runtime directory exists.
func EnsureRuntimeIntegrity() error {
	return os.MkdirAll(agentDataDir(), 0o755)
}

// ReadState loads persisted agent PID state.
func ReadState() (MonitorState, error) {
	data, err := os.ReadFile(stateFilePath())
	if err != nil {
		return MonitorState{}, err
	}
	var s MonitorState
	if err := json.Unmarshal(data, &s); err != nil {
		return MonitorState{}, err
	}
	return s, nil
}

// WriteState persists agent PID state.
func WriteState(s MonitorState) error {
	data, err := json.Marshal(s)
	if err != nil {
		return err
	}
	return os.WriteFile(stateFilePath(), data, 0o644)
}

// RemoveState deletes persisted agent PID state.
func RemoveState() error {
	err := os.Remove(stateFilePath())
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

// ProcessRunning reports whether pid is still alive on this OS.
func ProcessRunning(pid int) bool {
	if pid <= 0 {
		return false
	}
	return syscall.Kill(pid, 0) == nil
}

// LoadConfig reads optional monitor JSON config.
func LoadConfig(path string) (Config, error) {
	if strings.TrimSpace(path) == "" {
		return Config{}, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Config{}, nil
		}
		return Config{}, err
	}
	var c Config
	if err := json.Unmarshal(data, &c); err != nil {
		return Config{}, err
	}
	return c, nil
}

// GetRuntimeConfig loads persisted runtime JSON (hub URL, etc.).
func GetRuntimeConfig() (RuntimeConfig, error) {
	data, err := os.ReadFile(runtimeFilePath())
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return RuntimeConfig{}, nil
		}
		return RuntimeConfig{}, err
	}
	var c RuntimeConfig
	if err := json.Unmarshal(data, &c); err != nil {
		return RuntimeConfig{}, err
	}
	return c, nil
}

// APIServerAddress returns the local agent HTTP listen address.
func APIServerAddress() string {
	if v := strings.TrimSpace(os.Getenv("NEVARIX_AGENT_HTTP_ADDR")); v != "" {
		return v
	}
	return ":8080"
}

// APIToken returns the bearer token required for the local API.
// NEVARIX_AGENT_API_TOKEN takes precedence when set; otherwise the token is read from a
// dedicated secret file written by SaveAPITokenSecret (mode 0600).
func APIToken() string {
	if v := strings.TrimSpace(os.Getenv("NEVARIX_AGENT_API_TOKEN")); v != "" {
		return v
	}
	data, err := os.ReadFile(apiTokenSecretPath())
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

// SaveAPITokenSecret writes the API token to disk with owner-only permissions, replacing any prior value atomically.
func SaveAPITokenSecret(token string) error {
	token = strings.TrimSpace(token)
	if token == "" {
		return errors.New("API token is empty")
	}
	if err := EnsureRuntimeIntegrity(); err != nil {
		return err
	}
	dir := agentDataDir()
	f, err := os.CreateTemp(dir, ".api_token.")
	if err != nil {
		return err
	}
	tmpPath := f.Name()
	if _, err := f.WriteString(token); err != nil {
		_ = f.Close()
		_ = os.Remove(tmpPath)
		return err
	}
	if err := f.Sync(); err != nil {
		_ = f.Close()
		_ = os.Remove(tmpPath)
		return err
	}
	if err := f.Close(); err != nil {
		_ = os.Remove(tmpPath)
		return err
	}
	if err := os.Chmod(tmpPath, 0o600); err != nil {
		_ = os.Remove(tmpPath)
		return err
	}
	if err := os.Rename(tmpPath, apiTokenSecretPath()); err != nil {
		_ = os.Remove(tmpPath)
		return err
	}
	return nil
}
