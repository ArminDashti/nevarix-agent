package main

import (
	"fmt"
	"os"

	monitor "radar-agent/internal/domain/prober"
)

func main() {
	if err := ensureAgentRuntimeDir(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := monitor.EnsureRuntimeIntegrity(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := executeCommand(os.Args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
