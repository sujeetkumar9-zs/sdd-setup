package cmd

import (
	"fmt"
	"os/exec"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/sujeetkumar9-zs/sdd-setup/internal/mcp"
	"github.com/sujeetkumar9-zs/sdd-setup/internal/mempalace"
	"github.com/sujeetkumar9-zs/sdd-setup/internal/system"
	"github.com/sujeetkumar9-zs/sdd-setup/internal/templates"
)

var skipMine bool

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Complete one-time SDD setup for your project",
	Long: `Sets up everything needed for Spec-Driven Development:

  ✓ Detects project language
  ✓ Checks all prerequisites
  ✓ Creates Python virtual environment
  ✓ Installs language quality tools
  ✓ Initializes Mempalace knowledge graph
  ✓ Mines your codebase
  ✓ Configures Claude MCP servers
  ✓ Creates SKILL.md and slash commands
  ✓ Updates .gitignore

Usage:
  cd your-project
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

	lang := system.DetectLanguage()
	printDetectedLang(lang)

	steps := []step{
		{
			name:        "Checking prerequisites",
			description: "Python, Claude, mempalace" + langPrereqDesc(lang),
			fn:          func() error { return system.CheckPrerequisites(lang) },
		},
		{
			name:        "Creating Python virtual environment",
			description: "Isolated Python for mempalace",
			fn:          system.CreateVenv,
		},
		{
			name:        "Installing language tools",
			description: langToolDesc(lang),
			fn:          func() error { return system.InstallLanguageTools(lang) },
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

	// Verify everything is wired up correctly
	fmt.Printf("  %s Running post-setup verification...\n\n", color.BlueString("→"))
	if err := runVerify(cmd, args); err != nil {
		fmt.Println()
		color.Yellow("  Setup completed but some checks failed — see above for fix hints.")
		fmt.Println()
	}

	return nil
}

// ─────────────────────────────────────────────────────────
// Helpers
// ─────────────────────────────────────────────────────────

func langPrereqDesc(lang string) string {
	switch lang {
	case system.LangGo:
		return ", go"
	case system.LangNode:
		return ", node, npm"
	case system.LangRust:
		return ", cargo"
	default:
		return ""
	}
}

func langToolDesc(lang string) string {
	switch lang {
	case system.LangGo:
		return "golangci-lint, mockery"
	case system.LangNode:
		return "eslint"
	case system.LangPython:
		return "ruff, pytest"
	case system.LangRust:
		return "using rustup (no additional install)"
	case system.LangJava:
		return "using Maven/Gradle (no additional install)"
	default:
		return "none detected"
	}
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
	color.Blue("  SDD Setup")
	color.Blue("  Spec-Driven Development " + Version)
	color.Blue("═══════════════════════════════════════════")
	fmt.Println()
}

func printDetectedLang(lang string) {
	fmt.Printf("  %s Detected language: %s\n\n",
		color.BlueString("→"),
		color.CyanString(lang),
	)
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
		color.CyanString("/spec-to-pr <jira-key> [confluence-url ...]"))
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
