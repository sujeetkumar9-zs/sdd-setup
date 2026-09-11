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
var resumeSetup bool

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
	setupCmd.Flags().BoolVar(&skipMine, "skip-mine", false,
		"Skip the codebase mining step (faster setup)")
	setupCmd.Flags().BoolVar(&resumeSetup, "resume", false,
		"Resume a previously interrupted setup, skipping already-completed steps")
}

func runSetup(cmd *cobra.Command, args []string) error {
	printBanner()

	// Load or initialise state for resume support.
	var state *setupState
	if resumeSetup {
		s, err := loadState()
		if err != nil {
			return fmt.Errorf(
				"no interrupted setup found — run %s to start fresh\n  (looked for %s)",
				color.CyanString("sdd setup"), stateFile,
			)
		}
		state = s
		fmt.Printf("  %s Resuming previous setup (%d step(s) already completed)\n\n",
			color.YellowString("↻"), len(state.Completed))
	} else {
		state = &setupState{}
	}

	lang := system.DetectLanguage()
	printDetectedLang(lang)

	steps := []step{
		{
			id:          "check-prerequisites",
			name:        "Checking prerequisites",
			description: "Python, Claude, mempalace" + langPrereqDesc(lang),
			fn:          func() error { return system.CheckPrerequisites(lang) },
			rcaHint:     "One or more required tools are missing. Install each tool listed above, then re-run 'sdd setup'.",
		},
		{
			id:          "create-venv",
			name:        "Creating Python virtual environment",
			description: "Isolated Python for mempalace",
			fn:          system.CreateVenv,
			rcaHint:     "Check that python3 is installed ('python3 --version') and that you have write access to the current directory.",
		},
		{
			id:          "install-lang-tools",
			name:        "Installing language tools",
			description: langToolDesc(lang),
			fn:          func() error { return system.InstallLanguageTools(lang) },
			rcaHint:     "Check your network connection and that the package manager (go/npm/pip) is working correctly.",
		},
		{
			id:          "init-mempalace",
			name:        "Initializing Mempalace",
			description: "Creating .mempalace/palace/",
			fn:          mempalace.Init,
			rcaHint:     "Check that 'mempalace' is installed ('pipx install mempalace') and that you have write access to the current directory.",
		},
		{
			id:          "configure-mcp",
			name:        "Configuring Claude MCP servers",
			description: "Connecting mempalace to Claude",
			fn:          mcp.Configure,
			rcaHint:     "Check that 'claude' CLI is installed ('npm install -g @anthropic-ai/claude-code') and accessible on your PATH.",
		},
		{
			id:          "install-skills",
			name:        "Installing SDD skills and commands",
			description: "SKILL.md, slash commands, examples",
			fn:          templates.Install,
			rcaHint:     "Check that you have write access to the current directory and that '.claude/' can be created.",
		},
		{
			id:          "update-gitignore",
			name:        "Updating .gitignore",
			description: "Excluding .mempalace and .venv",
			fn:          system.UpdateGitignore,
			rcaHint:     "Check that you have write access to the .gitignore file in this directory.",
		},
	}

	// Add mining step unless skipped
	if !skipMine {
		mineStep := step{
			id:          "mine-codebase",
			name:        "Mining codebase",
			description: "Building knowledge graph (may take a few mins)",
			fn:          mempalace.Mine,
			rcaHint:     "Mining can fail if mempalace is not correctly installed or the .mempalace/palace directory is missing. Try 'sdd verify'.",
		}
		// Insert mining after init
		steps = append(steps[:4], append([]step{mineStep}, steps[4:]...)...)
	}

	// Execute all steps
	total := len(steps)
	for i, s := range steps {
		if state.isDone(s.id) {
			printStepSkipped(i+1, total, s.name)
			continue
		}

		printStep(i+1, total, s.name, s.description)

		if err := s.fn(); err != nil {
			printStepFail(s, err)
			// Save state so --resume can skip what succeeded so far.
			_ = state.save()
			return err
		}

		state.markDone(s.id)
		_ = state.save()
		printStepSuccess()
	}

	clearState()
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
	id          string
	name        string
	description string
	fn          func() error
	rcaHint     string
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

func printStepSkipped(current, total int, name string) {
	fmt.Printf(
		"\n  %s [%d/%d] %s\n  %s\n",
		color.HiBlackString("↷"),
		current, total,
		color.HiBlackString(name),
		color.HiBlackString("   already completed — skipping"),
	)
}

func printStepFail(s step, err error) {
	fmt.Println()
	color.Red("  ✗ Step failed: " + s.name)
	fmt.Println()
	color.Yellow("  ── Error ──────────────────────────────────")
	fmt.Printf("  %s\n", err.Error())
	fmt.Println()
	if s.rcaHint != "" {
		color.Yellow("  ── Root Cause Analysis ────────────────────")
		fmt.Printf("  %s\n", s.rcaHint)
		fmt.Println()
	}
	color.Yellow("  ── Next steps ─────────────────────────────")
	fmt.Printf("  • Fix the issue above, then resume:  %s\n", color.CyanString("sdd setup --resume"))
	fmt.Printf("  • Or restart from scratch:           %s\n", color.CyanString("sdd setup"))
	fmt.Printf("  • Diagnose your environment:         %s\n", color.CyanString("sdd verify"))
	fmt.Println()
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
