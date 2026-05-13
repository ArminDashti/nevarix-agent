package monitor

import (
	"errors"
	"os"
	"strings"
)

const agentRuntimeDirPath = "/home/.nevarix-server"

func agentRuntimeDir() string {
	return agentRuntimeDirPath
}

// EnsureRuntimeIntegrity verifies the agent can persist state under the runtime directory.
func EnsureRuntimeIntegrity() error {
	fi, err := os.Stat(agentRuntimeDirPath)
	if err != nil {
		return err
	}
	if !fi.IsDir() {
		return errors.New("agent runtime path is not a directory")
	}
	return nil
}

// APIServerAddress returns the HTTP listen address (host:port or :port).
func APIServerAddress() string {
	if a := strings.TrimSpace(os.Getenv("AGENT_HTTP_ADDR")); a != "" {
		return a
	}
	return ":8080"
}

// APIToken returns the bearer token required for the local API.
func APIToken() string {
	if t := strings.TrimSpace(os.Getenv("NEVARIX_AGENT_API_TOKEN")); t != "" {
		return t
	}
	return strings.TrimSpace(os.Getenv("AGENT_API_TOKEN"))
}
