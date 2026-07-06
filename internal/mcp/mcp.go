package mcp

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Configure sets up Claude MCP servers
func Configure() error {
	absPath, err := filepath.Abs(".mempalace/palace")
	if err != nil {
		return fmt.Errorf("could not resolve palace path: %w", err)
	}

	// Remove old broken entry if exists
	exec.Command("claude", "mcp", "remove", "mempalace").Run()

	// Add with correct command using absolute path so Claude can find
	// the palace regardless of which directory it is opened from.
	cmd := exec.Command(
		"claude", "mcp", "add", "mempalace",
		"--",
		"mempalace-mcp",
		"--palace", absPath,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
