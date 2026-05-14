package monitor

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestAPIToken_envOverridesFile(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("NEVARIX_AGENT_RUNTIME_DIR", dir)
	t.Setenv("NEVARIX_AGENT_API_TOKEN", "from-env")

	if err := SaveAPITokenSecret("from-file"); err != nil {
		t.Fatal(err)
	}
	if got := APIToken(); got != "from-env" {
		t.Fatalf("expected env token, got %q", got)
	}
}

func TestAPIToken_readsSecretFile(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("NEVARIX_AGENT_RUNTIME_DIR", dir)
	t.Setenv("NEVARIX_AGENT_API_TOKEN", "")

	if err := SaveAPITokenSecret("secret-value"); err != nil {
		t.Fatal(err)
	}
	if got := APIToken(); got != "secret-value" {
		t.Fatalf("expected file token, got %q", got)
	}
}

func TestSaveAPITokenSecret_filePermissions(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("file mode checks are unix-specific")
	}
	dir := t.TempDir()
	t.Setenv("NEVARIX_AGENT_RUNTIME_DIR", dir)

	if err := SaveAPITokenSecret("x"); err != nil {
		t.Fatal(err)
	}
	fi, err := os.Stat(filepath.Join(dir, "api_token.secret"))
	if err != nil {
		t.Fatal(err)
	}
	if fi.Mode().Perm()&0o777 != 0o600 {
		t.Fatalf("expected 0600, got %#o", fi.Mode().Perm()&0o777)
	}
}
