package main

import (
	"fmt"
	"os"
	"strings"

	monitor "nevarix-agent/internal/domain/prober"
)

func main() {
	skipRuntime := len(os.Args) >= 3 &&
		strings.EqualFold(strings.TrimSpace(os.Args[1]), "connect") &&
		strings.EqualFold(strings.TrimSpace(os.Args[2]), "hub")

	if !skipRuntime {
		if err := monitor.EnsureRuntimeIntegrity(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	if err := executeCommand(os.Args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
