package system

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// CheckPrerequisites verifies all required tools are installed
func CheckPrerequisites() error {
	required := map[string]string{
		"go":        "Install from: https://go.dev/dl/",
		"python3":   "Run: brew install python",
		"claude":    "Run: npm install -g @anthropic-ai/claude-code",
		"mempalace": "Run: pipx install mempalace",
		"docker":    "Install from: https://docker.com",
	}

	var missing []string
	for tool, hint := range required {
		if _, err := exec.LookPath(tool); err != nil {
			missing = append(missing,
				fmt.Sprintf("%s (%s)", tool, hint))
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf(
			"missing tools:\n  • %s",
			strings.Join(missing, "\n  • "),
		)
	}
	return nil
}

// CreateVenv creates Python virtual environment
func CreateVenv() error {
	if _, err := os.Stat(".venv"); os.IsNotExist(err) {
		return exec.Command("python3", "-m", "venv", ".venv").Run()
	}
	return nil
}

// InstallGoTools installs required Go development tools
func InstallGoTools() error {
	tools := []string{
		"github.com/golangci/golangci-lint/cmd/golangci-lint@latest",
		"github.com/vektra/mockery/v2@latest",
		"golang.org/x/tools/cmd/goimports@latest",
		"gotest.tools/gotestsum@latest",
	}

	for _, tool := range tools {
		cmd := exec.Command("go", "install", tool)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to install %s: %w", tool, err)
		}
	}
	return nil
}

// UpdateGitignore adds required entries to .gitignore
func UpdateGitignore() error {
	content, _ := os.ReadFile(".gitignore")
	existing := string(content)

	additions := "\n# SDD - Spec Driven Development\n"
	changed := false

	entries := []string{".mempalace/", ".venv/", "__pycache__/"}
	for _, entry := range entries {
		if !strings.Contains(existing, entry) {
			additions += entry + "\n"
			changed = true
		}
	}

	if changed {
		f, err := os.OpenFile(".gitignore",
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = f.WriteString(additions)
		return err
	}
	return nil
}
