package mempalace

import (
	"fmt"
	"os"
	"os/exec"
)

const palacePath = "./.mempalace/palace"

// Init initializes mempalace in the current directory
func Init() error {
	if _, err := os.Stat(".mempalace"); os.IsNotExist(err) {
		return run("mempalace", "init", ".")
	}
	return nil
}

// Mine runs mempalace mining to build knowledge graph
func Mine() error {
	return run("mempalace", "--palace", palacePath, "mine", "./")
}

// Status shows current palace status
func Status() error {
	return run("mempalace", "--palace", palacePath, "status")
}

// Search queries the knowledge graph
func Search(query string) error {
	return run("mempalace", "--palace", palacePath, "search", query)
}

// WakeUp shows wake-up context
func WakeUp() error {
	return run("mempalace", "--palace", palacePath, "wake-up")
}

// GetMCPCommand returns the correct MCP setup command
func GetMCPCommand() (string, error) {
	out, err := exec.Command(
		"mempalace", "--palace", palacePath, "mcp",
	).Output()
	if err != nil {
		return "", fmt.Errorf("failed to get MCP command: %w", err)
	}
	return string(out), nil
}

func run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
