package monitor

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"syscall"
)

const stateFileName = "agent.state"

func statePath() string {
	return filepath.Join(agentRuntimeDir(), stateFileName)
}

// ReadState loads persisted agent PID state.
func ReadState() (MonitorState, error) {
	var s MonitorState
	b, err := os.ReadFile(statePath())
	if err != nil {
		return s, err
	}
	if err := json.Unmarshal(b, &s); err != nil {
		return s, err
	}
	return s, nil
}

// WriteState persists agent PID state.
func WriteState(s MonitorState) error {
	b, err := json.Marshal(s)
	if err != nil {
		return err
	}
	return os.WriteFile(statePath(), b, 0o600)
}

// RemoveState deletes persisted agent state.
func RemoveState() error {
	err := os.Remove(statePath())
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}

// ProcessRunning reports whether pid is an alive process on this host.
func ProcessRunning(pid int) bool {
	if pid <= 0 {
		return false
	}
	p, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	err = p.Signal(syscall.Signal(0))
	return err == nil
}
