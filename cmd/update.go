package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/sujeetkumar9-zs/sdd-setup/internal/templates"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update SDD skills and commands to the latest version",
	Long: `Overwrites all SDD skill and command files with the versions bundled in this
binary, without touching .mempalace/, .venv/, or any project files.

Run this after upgrading the sdd binary to pick up new or improved skills.`,
	RunE: runUpdate,
}

func runUpdate(_ *cobra.Command, _ []string) error {
	// Guard: must be a project that has already been set up.
	if _, err := os.Stat(".mempalace"); os.IsNotExist(err) {
		return fmt.Errorf(
			"no .mempalace directory found — run %s first",
			color.CyanString("sdd setup"),
		)
	}

	fmt.Println()
	fmt.Printf("  %s Updating SDD skills and commands...\n\n",
		color.BlueString("→"))

	n, err := templates.InstallForce()
	if err != nil {
		fmt.Printf("  %s Update failed: %s\n\n",
			color.RedString("✗"), err)
		return err
	}

	fmt.Printf("  %s %d skill/command files updated\n",
		color.GreenString("✓"), n)
	fmt.Printf("  %s .claude/CLAUDE.md refreshed\n\n",
		color.GreenString("✓"))

	color.Green("  Skills are up to date. No other project files were changed.")
	fmt.Println()
	return nil
}
