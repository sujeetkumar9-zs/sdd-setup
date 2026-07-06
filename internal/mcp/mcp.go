package mcp

import (
	"os/exec"
)

// Configure sets up Claude MCP servers
func Configure() error {
	// Remove old broken entry if exists
	exec.Command("claude", "mcp", "remove", "mempalace").Run()

	// Add with correct command
	return exec.Command(
		"claude", "mcp", "add", "mempalace",
		"--",
		"mempalace-mcp",
		"--palace", ".mempalace/palace",
	).Run()
}
