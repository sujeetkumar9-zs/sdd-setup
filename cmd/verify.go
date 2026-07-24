package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verify SDD setup is working correctly",
	Long: `Runs a comprehensive check of your SDD setup:

  ✓ Mempalace installed and working
  ✓ Palace initialized with data
  ✓ MCP server registered and connected
  ✓ Claude Code installed
  ✓ Language quality tools available
  ✓ Virtual environment exists
  ✓ .gitignore configured correctly

Run this if something seems wrong with your setup.`,
	RunE: runVerify,
}

type check struct {
	name    string
	fn      func() error
	fixHint string
}

func runVerify(cmd *cobra.Command, args []string) error {
	fmt.Println()
	fmt.Printf("  %s Verifying SDD setup...\n\n",
		color.BlueString("→"))

	checks := []check{
		{
			name:    "Mempalace installed",
			fn:      checkMempalaceInstalled,
			fixHint: "Run: pipx install mempalace",
		},
		{
			name:    ".mempalace directory exists",
			fn:      checkMempalaceDir,
			fixHint: "Run: sdd setup",
		},
		{
			name:    "Palace has indexed data",
			fn:      checkPalaceData,
			fixHint: "Run: sdd mine",
		},
		{
			name:    "Mempalace MCP registered",
			fn:      checkMCPRegistered,
			fixHint: "Run: claude mcp add mempalace -- mempalace-mcp --palace .mempalace/palace",
		},
		{
			name:    "Claude Code installed",
			fn:      checkClaudeInstalled,
			fixHint: "Run: npm install -g @anthropic-ai/claude-code",
		},
		{
			name:    "Python venv exists",
			fn:      checkVenv,
			fixHint: "Run: python3 -m venv .venv",
		},
		{
			name:    "SKILL.md exists",
			fn:      checkSkillMD,
			fixHint: "Run: sdd setup",
		},
		{
			name:    ".gitignore configured",
			fn:      checkGitignore,
			fixHint: "Run: sdd setup",
		},
	}

	// Add Go-specific tool checks only for Go projects
	if isGoProject() {
		checks = append(checks,
			check{
				name:    "golangci-lint installed",
				fn:      func() error { return checkTool("golangci-lint") },
				fixHint: "Run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest",
			},
			check{
				name:    "mockery installed",
				fn:      func() error { return checkTool("mockery") },
				fixHint: "Run: go install github.com/vektra/mockery/v2@latest",
			},
		)
	}

	passed := 0
	failed := 0
	hints := []string{}

	for _, c := range checks {
		err := c.fn()
		if err != nil {
			fmt.Printf("  %s %s\n",
				color.RedString("✗"),
				c.name)
			hints = append(hints, c.fixHint)
			failed++
		} else {
			fmt.Printf("  %s %s\n",
				color.GreenString("✓"),
				c.name)
			passed++
		}
	}

	fmt.Println()
	fmt.Printf("  Results: %s passed, %s failed\n",
		color.GreenString("%d", passed),
		color.RedString("%d", failed))
	fmt.Println()

	if failed > 0 {
		color.Yellow("  Fix hints:")
		for _, hint := range hints {
			fmt.Printf("  • %s\n", hint)
		}
		fmt.Println()
		return fmt.Errorf("%d checks failed", failed)
	}

	color.Green("  ✅ All checks passed! SDD is ready to use.")
	fmt.Println()
	return nil
}

func isGoProject() bool {
	_, err := os.Stat("go.mod")
	return err == nil
}

func checkMempalaceInstalled() error {
	_, err := exec.LookPath("mempalace")
	return err
}

func checkMempalaceDir() error {
	if _, err := os.Stat(".mempalace/palace"); os.IsNotExist(err) {
		return fmt.Errorf("not found")
	}
	return nil
}

func checkPalaceData() error {
	// mempalace 3.x uses ChromaDB (chroma.sqlite3); older versions used index.db
	for _, candidate := range []string{
		".mempalace/palace/chroma.sqlite3",
		".mempalace/palace/index.db",
	} {
		info, err := os.Stat(candidate)
		if err == nil && info.Size() > 0 {
			return nil
		}
	}
	return fmt.Errorf("palace has no data — run: sdd mine")
}

func checkMCPRegistered() error {
	out, err := exec.Command("claude", "mcp", "list").Output()
	if err != nil {
		return fmt.Errorf("claude mcp list failed")
	}
	if !strings.Contains(string(out), "mempalace") {
		return fmt.Errorf("mempalace not in MCP list")
	}
	return nil
}

func checkClaudeInstalled() error {
	_, err := exec.LookPath("claude")
	return err
}

func checkTool(name string) error {
	_, err := exec.LookPath(name)
	return err
}

func checkVenv() error {
	if _, err := os.Stat(".venv"); os.IsNotExist(err) {
		return fmt.Errorf("not found")
	}
	return nil
}

func checkSkillMD() error {
	path := ".claude/skills/spec-to-pr/SKILL.md"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("not found")
	}
	return nil
}

func checkGitignore() error {
	data, err := os.ReadFile(".gitignore")
	if err != nil {
		return fmt.Errorf("no .gitignore found")
	}
	content := string(data)
	if !strings.Contains(content, ".mempalace") {
		return fmt.Errorf(".mempalace not in .gitignore")
	}
	return nil
}
