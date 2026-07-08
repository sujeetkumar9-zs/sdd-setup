package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var teardownCmd = &cobra.Command{
	Use:   "teardown",
	Short: "Remove SDD setup from this project",
	Long: `Removes all SDD setup from the current project:

  • Deregisters the mempalace MCP server from Claude
  • Deletes .mempalace/ (knowledge graph)
  • Deletes .venv/ (Python virtual environment)

Does not remove .claude/commands, SKILL.md, or .golangci.yml.
Run sdd setup again to restore.`,
	RunE: runTeardown,
}

func runTeardown(cmd *cobra.Command, args []string) error {
	fmt.Println()
	color.Yellow("  This will remove .mempalace/ and .venv/ from this project.")
	fmt.Printf("  Type %s to confirm: ", color.RedString("yes"))

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	if strings.TrimSpace(scanner.Text()) != "yes" {
		fmt.Println("  Aborted.")
		fmt.Println()
		return nil
	}

	fmt.Println()

	steps := []struct {
		desc string
		fn   func() error
	}{
		{
			desc: "Deregistering mempalace MCP",
			fn: func() error {
				exec.Command("claude", "mcp", "remove", "mempalace").Run()
				return nil
			},
		},
		{
			desc: "Removing .mempalace/",
			fn:   func() error { return os.RemoveAll(".mempalace") },
		},
		{
			desc: "Removing .venv/",
			fn:   func() error { return os.RemoveAll(".venv") },
		},
	}

	for _, s := range steps {
		fmt.Printf("  %s %s...\n", color.BlueString("→"), s.desc)
		if err := s.fn(); err != nil {
			fmt.Printf("  %s %s: %s\n", color.RedString("✗"), s.desc, err.Error())
			return err
		}
		fmt.Printf("  %s Done\n", color.GreenString("✓"))
	}

	fmt.Println()
	color.Green("  Teardown complete. Run 'sdd setup' to set up again.")
	fmt.Println()
	return nil
}
