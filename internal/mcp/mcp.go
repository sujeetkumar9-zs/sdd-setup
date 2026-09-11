package mcp

import (
	"fmt"
	"os/exec"
	"path/filepath"
)

// Configure sets up Claude MCP servers
func Configure() error {
	if _, err := exec.LookPath("claude"); err != nil {
		return fmt.Errorf("claude CLI not found in PATH\n  RCA: 'claude' binary is not installed or not on your PATH\n  Fix: npm install -g @anthropic-ai/claude-code")
	}

	absPath, err := filepath.Abs(".mempalace/palace")
	if err != nil {
		return fmt.Errorf("could not resolve palace path: %w", err)
	}

	// Remove old broken entry if exists — ignore errors (may not exist yet)
	exec.Command("claude", "mcp", "remove", "mempalace").Run() //nolint:errcheck

	// Add with correct command using absolute path so Claude can find
	// the palace regardless of which directory it is opened from.
	out, err := exec.Command(
		"claude", "mcp", "add", "mempalace",
		"--",
		"mempalace-mcp",
		"--palace", absPath,
	).CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to register mempalace MCP server\n  RCA: %w\n  Output: %s", err, string(out))
	}
	return nil
}
