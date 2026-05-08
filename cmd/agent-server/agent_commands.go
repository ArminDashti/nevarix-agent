package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"syscall"
	"time"

	posttohub "radar-agent/internal/domain/post_to_hub"
	monitor "radar-agent/internal/domain/prober"
	"radar-agent/internal/http/api"
)

// startCommand starts the agent process.
func startCommand() error {
	if state, err := monitor.ReadState(); err == nil && monitor.ProcessRunning(state.PID) {
		fmt.Printf("agent is already running (pid %d)\n", state.PID)
		return nil
	}
	_ = monitor.RemoveState()

	executable, err := os.Executable()
	if err != nil {
		return err
	}
	cmd := exec.Command(executable, "agent")
	logFile, err := os.OpenFile("agent.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	if err := cmd.Start(); err != nil {
		_ = logFile.Close()
		return err
	}
	pid := cmd.Process.Pid
	_ = logFile.Close()
	if err := cmd.Process.Release(); err != nil {
		return err
	}
	if err := monitor.WriteState(monitor.MonitorState{PID: pid, StartedAtUnix: time.Now().Unix()}); err != nil {
		return err
	}
	fmt.Printf("agent started (pid %d)\n", pid)
	return nil
}

func stopCommand() error {
	state, err := monitor.ReadState()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("agent is not running")
			return nil
		}
		return err
	}
	process, err := os.FindProcess(state.PID)
	if err != nil {
		return err
	}
	if err := process.Signal(syscall.SIGTERM); err != nil {
		return err
	}
	if err := monitor.RemoveState(); err != nil {
		return err
	}
	fmt.Printf("agent stopped (pid %d)\n", state.PID)
	return nil
}

func runRestart() error {
	if err := stopCommand(); err != nil {
		return err
	}
	return startCommand()
}

func runStatus() error {
	state, err := monitor.ReadState()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("agent is not running")
			return nil
		}
		return err
	}
	if !monitor.ProcessRunning(state.PID) {
		fmt.Printf("agent state exists but pid %d is not running\n", state.PID)
		return nil
	}
	fmt.Printf("agent is running (pid %d)\n", state.PID)
	return nil
}

func startAPIServer() error {
	addr := monitor.APIServerAddress()
	token := monitor.APIToken()

	if token == "" {
		return fmt.Errorf("API token not configured")
	}

	server := &http.Server{
		Addr:    addr,
		Handler: api.NewRouter(token),
	}

	fmt.Printf("Starting API server on %s\n", addr)
	log.Printf("API server listening on %s", addr)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start API server: %w", err)
	}

	return nil
}

// runMonitor runs the monitoring loop that collects data and stores it in cache
func runMonitor() error {
	// Load the monitor config
	cfg, err := monitor.LoadConfig(monitor.ConfigFile)
	if err != nil {
		return fmt.Errorf("failed to load config: %v", err)
	}

	runtimeCfg, err := monitor.GetRuntimeConfig()
	if err != nil {
		return fmt.Errorf("failed to get runtime config: %v", err)
	}

	if runtimeCfg.API.HubBaseURL == "" {
		return fmt.Errorf("hub base URL not configured")
	}

	token := monitor.APIToken()
	if token == "" {
		return fmt.Errorf("API token not configured")
	}

	fmt.Println("Starting monitoring loop...")

	for {
		checkedAt := time.Now()
		if err := posttohub.CollectAndStoreInCache(runtimeCfg.API.HubBaseURL, token, checkedAt); err != nil {
			fmt.Printf("Error collecting and storing data: %v\n", err)
			// Continue monitoring despite errors
		}

		// Sleep for the configured interval
		interval := 60 * time.Second // Default 1 minute
		if cfg.IntervalSeconds > 0 {
			interval = time.Duration(cfg.IntervalSeconds) * time.Second
		}
		time.Sleep(interval)
	}
}
