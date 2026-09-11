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

var teardownFull bool

var teardownCmd = &cobra.Command{
	Use:   "teardown",
	Short: "Remove SDD setup from this project",
	Long: `Removes SDD setup from the current project:

  • Deregisters the mempalace MCP server from Claude
  • Deletes .mempalace/ (knowledge graph)
  • Deletes .venv/ (Python virtual environment)

With --full, also removes:
  • .claude/commands/ (SDD slash commands)
  • SKILL.md (SDD skills reference)

Run sdd setup again to restore.`,
	RunE: runTeardown,
}

func init() {
	teardownCmd.Flags().BoolVar(
		&teardownFull,
		"full",
		false,
		"Also remove .claude/commands/ and SKILL.md",
	)
}

func runTeardown(cmd *cobra.Command, args []string) error {
	fmt.Println()
	if teardownFull {
		color.Yellow("  This will remove .mempalace/, .venv/, .claude/commands/, and SKILL.md from this project.")
	} else {
		color.Yellow("  This will remove .mempalace/ and .venv/ from this project.")
		fmt.Printf("  Use %s to also remove .claude/commands/ and SKILL.md\n", color.CyanString("--full"))
	}
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
				exec.Command("claude", "mcp", "remove", "mempalace").Run() //nolint:errcheck
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

	if teardownFull {
		steps = append(steps,
			struct {
				desc string
				fn   func() error
			}{
				desc: "Removing .claude/commands/",
				fn:   func() error { return os.RemoveAll(".claude/commands") },
			},
			struct {
				desc string
				fn   func() error
			}{
				desc: "Removing SKILL.md",
				fn:   func() error { return os.Remove("SKILL.md") },
			},
		)
	}

	for _, s := range steps {
		fmt.Printf("  %s %s...\n", color.BlueString("→"), s.desc)
		if err := s.fn(); err != nil && !os.IsNotExist(err) {
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
