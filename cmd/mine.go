package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/sujeetkumar9-zs/sdd-setup/internal/mempalace"
	"github.com/spf13/cobra"
)

var mineCmd = &cobra.Command{
	Use:   "mine",
	Short: "Update codebase knowledge graph after coding",
	Long: `Re-mines your codebase to update Mempalace knowledge graph.

Run this after:
  - Implementing a new feature
  - Adding new packages or types
  - Significant code changes

This ensures Claude has up-to-date knowledge
of your codebase for future features.`,
	RunE: runMine,
}

func runMine(cmd *cobra.Command, args []string) error {
	fmt.Println()
	fmt.Printf("  %s Updating knowledge graph...\n",
		color.BlueString("→"))
	fmt.Println()

	if err := mempalace.Mine(); err != nil {
		return fmt.Errorf("mining failed: %w", err)
	}

	if err := mempalace.Status(); err != nil {
		return err
	}

	fmt.Println()
	fmt.Printf("  %s Knowledge graph updated!\n",
		color.GreenString("✓"))
	fmt.Printf("  %s Claude now knows about your new code\n",
		color.GreenString("✓"))
	fmt.Println()

	return nil
}
