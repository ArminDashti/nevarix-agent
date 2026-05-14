package main

import (
	"errors"
	"fmt"
	"strings"
)

var errUnknownCommand = errors.New("unknown command")

func executeCommand(args []string) error {
	if len(args) >= 2 {
		first := strings.ToLower(strings.TrimSpace(args[0]))
		second := strings.ToLower(strings.TrimSpace(args[1]))
		if first == "connect" && second == "hub" {
			return runConnectHub(args[2:])
		}
	}

	command := "start"
	if len(args) > 0 {
		command = strings.ToLower(strings.TrimSpace(args[0]))
	}

	switch command {
	case "", "start":
		return startCommand()
	case "agent":
		return startAPIServer()
	case "monitor":
		return runMonitor()
	case "send-to-hub":
		return runSendToHub()
	case "stop":
		return stopCommand()
	case "restart":
		return runRestart()
	case "status":
		return runStatus()
	case "version":
		fmt.Printf("agent version %s\n", buildVersion)
		return nil
	default:
		return errUnknownCommand
	}
}
