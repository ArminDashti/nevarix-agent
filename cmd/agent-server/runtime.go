package main

import "os"

func ensureAgentRuntimeDir() error {
	return os.MkdirAll("/home/.nevarix-server", 0o755) // 0o755 is the permission for the directory. readable and writable by the owner and the group.
}
