package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type hubAuthPayload struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	SecretKey string `json:"secret_key"`
}

type connectHubFlags struct {
	address   string
	username  string
	password  string
	secretKey string
}

func runConnectHub(flagArgs []string) error {
	flags, err := parseConnectHubFlags(flagArgs)
	if err != nil {
		return err
	}
	if flags.address == "" {
		return fmt.Errorf("--address is required")
	}
	if flags.username == "" {
		return fmt.Errorf("--username is required")
	}
	if flags.password == "" {
		return fmt.Errorf("--password is required")
	}
	if flags.secretKey == "" {
		return fmt.Errorf("--secret-key is required")
	}

	base := strings.TrimRight(strings.TrimSpace(flags.address), "/")
	authURL := base + "/api/v1/auth"

	body, err := json.Marshal(hubAuthPayload{
		Username:  flags.username,
		Password:  flags.password,
		SecretKey: flags.secretKey,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, authURL, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return err
	}

	fmt.Printf("POST %s -> %s\n", authURL, resp.Status)
	if len(respBody) > 0 {
		fmt.Println(string(respBody))
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("hub auth request failed with status %s", resp.Status)
	}
	return nil
}

func parseConnectHubFlags(args []string) (connectHubFlags, error) {
	var f connectHubFlags
	for _, raw := range args {
		if !strings.HasPrefix(raw, "--") {
			return f, fmt.Errorf("unexpected argument %q (expected --name=value)", raw)
		}
		raw = strings.TrimPrefix(raw, "--")
		key, val, ok := strings.Cut(raw, "=")
		if !ok {
			return f, fmt.Errorf("flag %q must use --name=value form", "--"+raw)
		}
		key = strings.ToLower(strings.TrimSpace(key))
		val = strings.TrimSpace(val)
		switch key {
		case "address":
			f.address = val
		case "username":
			f.username = val
		case "password":
			f.password = val
		case "secret-key":
			f.secretKey = val
		default:
			return f, fmt.Errorf("unknown flag --%s", key)
		}
	}
	return f, nil
}
