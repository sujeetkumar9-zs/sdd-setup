package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/fatih/color"
	"github.com/sujeetkumar9-zs/sdd-setup/internal/mcp"
	"github.com/sujeetkumar9-zs/sdd-setup/internal/mempalace"
	"github.com/sujeetkumar9-zs/sdd-setup/internal/system"
	"github.com/sujeetkumar9-zs/sdd-setup/internal/templates"
	"github.com/spf13/cobra"
)

var skipMine bool

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Complete one-time SDD setup for your Go project",
	Long: `Sets up everything needed for Spec-Driven Development:

  ✓ Checks all prerequisites
  ✓ Creates Python virtual environment
  ✓ Installs Go quality tools
  ✓ Initializes Mempalace knowledge graph
  ✓ Mines your codebase
  ✓ Configures Claude MCP servers
  ✓ Creates SKILL.md and slash commands
  ✓ Updates .gitignore

Usage:
  cd your-go-project
  sdd setup`,
	RunE: runSetup,
}

func init() {
	setupCmd.Flags().BoolVar(
		&skipMine,
		"skip-mine",
		false,
		"Skip the codebase mining step (faster setup)",
	)
}

func runSetup(cmd *cobra.Command, args []string) error {
	printBanner()

	// Verify Go project
	if _, err := os.Stat("go.mod"); os.IsNotExist(err) {
		return fmt.Errorf(
			"no go.mod found\n" +
				"  → Run this command from your Go project root",
		)
	}

	// Define all setup steps
	steps := []step{
		{
			name:        "Checking prerequisites",
			description: "Go, Python, Claude, Docker",
			fn:          system.CheckPrerequisites,
		},
		{
			name:        "Creating Python virtual environment",
			description: "Isolated Python for mempalace",
			fn:          system.CreateVenv,
		},
		{
			name:        "Installing Go quality tools",
			description: "golangci-lint, mockery, goimports",
			fn:          system.InstallGoTools,
		},
		{
			name:        "Initializing Mempalace",
			description: "Creating .mempalace/palace/",
			fn:          mempalace.Init,
		},
		{
			name:        "Configuring Claude MCP servers",
			description: "Connecting mempalace to Claude",
			fn:          mcp.Configure,
		},
		{
			name:        "Installing SDD skills and commands",
			description: "SKILL.md, slash commands, examples",
			fn:          templates.Install,
		},
		{
			name:        "Updating .gitignore",
			description: "Excluding .mempalace and .venv",
			fn:          system.UpdateGitignore,
		},
	}

	// Add mining step unless skipped
	if !skipMine {
		mineStep := step{
			name:        "Mining codebase",
			description: "Building knowledge graph (may take a few mins)",
			fn:          mempalace.Mine,
		}
		// Insert mining after init
		steps = append(steps[:4], append([]step{mineStep}, steps[4:]...)...)
	}

	// Execute all steps
	total := len(steps)
	for i, s := range steps {
		printStep(i+1, total, s.name, s.description)

		if err := s.fn(); err != nil {
			printStepFail(s.name, err)
			return err
		}

		printStepSuccess()
	}

	printSetupComplete()
	return nil
}

// ─────────────────────────────────────────────────────────
// Types
// ─────────────────────────────────────────────────────────

type step struct {
	name        string
	description string
	fn          func() error
}

// ─────────────────────────────────────────────────────────
// Print helpers
// ─────────────────────────────────────────────────────────

func printBanner() {
	fmt.Println()
	color.Blue("═══════════════════════════════════════════")
	color.Blue("  SDD Setup - Kroger Go Projects")
	color.Blue("  Spec-Driven Development " + Version)
	color.Blue("═══════════════════════════════════════════")
	fmt.Println()
}

func printStep(current, total int, name, desc string) {
	fmt.Printf(
		"\n  %s [%d/%d] %s\n  %s %s\n",
		color.BlueString("→"),
		current, total,
		color.WhiteString(name),
		color.HiBlackString("   "),
		color.HiBlackString(desc),
	)
}

func printStepSuccess() {
	fmt.Printf("  %s Done\n", color.GreenString("✓"))
}

func printStepFail(name string, err error) {
	fmt.Printf("  %s Failed: %s\n",
		color.RedString("✗"),
		color.RedString(err.Error()),
	)
	fmt.Println()
	color.Yellow("  Troubleshooting:")
	fmt.Printf("  Run %s for help\n",
		color.CyanString("sdd verify"),
	)
}

func printSetupComplete() {
	fmt.Println()
	color.Green("═══════════════════════════════════════════")
	color.Green("  ✅ SDD SETUP COMPLETE!")
	color.Green("═══════════════════════════════════════════")
	fmt.Println()
	fmt.Println("  Next steps:")
	fmt.Println()
	fmt.Printf("  1. Activate venv:    %s\n",
		color.CyanString("source .venv/bin/activate"))
	fmt.Printf("  2. Start Claude:     %s\n",
		color.CyanString("claude"))
	fmt.Printf("  3. Start a feature:  %s\n",
		color.CyanString("/spec-to-pr <jira-url>"))
	fmt.Printf("  4. After coding:     %s\n",
		color.CyanString("sdd mine"))
	fmt.Printf("  5. Before PR:        %s\n",
		color.CyanString("sdd quality"))
	fmt.Println()
	fmt.Printf("  Need help? Run %s\n",
		color.CyanString("sdd --help"),
	)
	fmt.Println()
}

// Ensure exec is available for setup
var _ = exec.Command
