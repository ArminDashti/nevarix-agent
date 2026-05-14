package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	monitor "nevarix-agent/internal/domain/prober"
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

// hubTokenJSONKeys lists JSON fields that may carry the local agent API token from the hub auth response.
var hubTokenJSONKeys = []string{
	"api_token", "access_token", "token", "bearer_token", "agent_api_token",
}

// extractTokenFromHubAuthResponse pulls a bearer token from common hub response shapes.
func extractTokenFromHubAuthResponse(body []byte) string {
	var root map[string]any
	if err := json.Unmarshal(body, &root); err != nil {
		return ""
	}
	if t := tokenFromJSONObject(root); t != "" {
		return t
	}
	if raw, ok := root["data"]; ok {
		switch nested := raw.(type) {
		case map[string]any:
			return tokenFromJSONObject(nested)
		}
	}
	return ""
}

// tokenFromJSONObject returns the first non-empty string among known token field names.
func tokenFromJSONObject(obj map[string]any) string {
	for _, key := range hubTokenJSONKeys {
		raw, ok := obj[key]
		if !ok {
			continue
		}
		s, ok := raw.(string)
		if ok && strings.TrimSpace(s) != "" {
			return strings.TrimSpace(s)
		}
	}
	return ""
}

// redactSensitiveHubFields replaces secret values so successful auth output can be logged safely.
func redactSensitiveHubFields(obj map[string]any) {
	sensitive := []string{
		"token", "access_token", "api_token", "bearer_token", "agent_api_token",
		"password", "secret_key", "refresh_token",
	}
	for _, k := range sensitive {
		if _, ok := obj[k]; ok {
			obj[k] = "<redacted>"
		}
	}
	if raw, ok := obj["data"]; ok {
		if nested, ok := raw.(map[string]any); ok {
			redactSensitiveHubFields(nested)
		}
	}
}

// redactHubAuthResponseForDisplay returns JSON text safe to print after hub authentication.
func redactHubAuthResponseForDisplay(body []byte) string {
	var root map[string]any
	if err := json.Unmarshal(body, &root); err != nil {
		return strings.TrimSpace(string(body))
	}
	redactSensitiveHubFields(root)
	out, err := json.MarshalIndent(root, "", "  ")
	if err != nil {
		return strings.TrimSpace(string(body))
	}
	return string(out)
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
		fmt.Println(redactHubAuthResponseForDisplay(respBody))
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("hub auth request failed with status %s", resp.Status)
	}

	if t := extractTokenFromHubAuthResponse(respBody); t != "" {
		if err := monitor.SaveAPITokenSecret(t); err != nil {
			return fmt.Errorf("hub auth succeeded but could not store API token securely: %w", err)
		}
		fmt.Println("Agent API token stored with owner-only file permissions under the agent data directory.")
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
