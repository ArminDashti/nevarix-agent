package main

import (
	"strings"
	"testing"
)

func TestExtractTokenFromHubAuthResponse(t *testing.T) {
	cases := []struct {
		name string
		json string
		want string
	}{
		{"top_level_token", `{"token":"abc"}`, "abc"},
		{"nested_data", `{"data":{"access_token":"xyz"}}`, "xyz"},
		{"invalid_json", `{`, ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := extractTokenFromHubAuthResponse([]byte(tc.json)); got != tc.want {
				t.Fatalf("got %q want %q", got, tc.want)
			}
		})
	}
}

func TestRedactHubAuthResponseForDisplay(t *testing.T) {
	in := `{"token":"secret","ok":true}`
	out := redactHubAuthResponseForDisplay([]byte(in))
	if !strings.Contains(out, "redacted") || !strings.Contains(out, "ok") {
		t.Fatalf("unexpected output: %s", out)
	}
}
