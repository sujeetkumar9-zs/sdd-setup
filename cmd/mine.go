package cmd

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/sujeetkumar9-zs/sdd-setup/internal/mempalace"
	"github.com/spf13/cobra"
)

var mineChanged bool

var mineCmd = &cobra.Command{
	Use:   "mine",
	Short: "Update codebase knowledge graph after coding",
	Long: `Re-mines your codebase to update Mempalace knowledge graph.

Run this after:
  - Implementing a new feature
  - Adding new packages or types
  - Significant code changes

Use --changed to mine only files git reports as modified (faster for large repos).

This ensures Claude has up-to-date knowledge
of your codebase for future features.`,
	RunE: runMine,
}

func init() {
	mineCmd.Flags().BoolVar(&mineChanged, "changed", false,
		"Mine only directories with git-modified files (staged, unstaged, and last commit)")
}

func runMine(cmd *cobra.Command, args []string) error {
	fmt.Println()

	if mineChanged {
		return runMineChanged()
	}

	fmt.Printf("  %s Updating knowledge graph...\n", color.BlueString("→"))
	fmt.Println()

	if err := mempalace.Mine(); err != nil {
		return fmt.Errorf("mining failed: %w", err)
	}

	return printMineSuccess()
}

func runMineChanged() error {
	dirs, err := gitChangedDirs()
	if err != nil {
		return fmt.Errorf("could not determine changed files: %w", err)
	}

	if len(dirs) == 0 {
		fmt.Printf("  %s No changed files detected — running full mine instead\n\n",
			color.YellowString("!"))
		if err := mempalace.Mine(); err != nil {
			return fmt.Errorf("mining failed: %w", err)
		}
		return printMineSuccess()
	}

	fmt.Printf("  %s Mining %d changed director%s...\n",
		color.BlueString("→"), len(dirs), pluralY(len(dirs)))
	for _, d := range dirs {
		fmt.Printf("      %s\n", d)
	}
	fmt.Println()

	if err := mempalace.MineChanged(dirs); err != nil {
		return fmt.Errorf("mining failed: %w", err)
	}

	return printMineSuccess()
}

// gitChangedDirs returns unique parent directories of all files that git
// considers modified: staged, unstaged, and changed in the last commit.
func gitChangedDirs() ([]string, error) {
	seen := map[string]struct{}{}
	var dirs []string

	gitSets := [][]string{
		{"diff", "--name-only", "HEAD"},     // unstaged
		{"diff", "--name-only", "--cached"}, // staged
		{"diff", "--name-only", "HEAD~1"},   // last commit (best-effort; may fail on first commit)
	}

	for _, gitArgs := range gitSets {
		out, _ := exec.Command("git", gitArgs...).Output() // ignore error — e.g. no prior commit
		for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
			if line == "" {
				continue
			}
			dir := filepath.Dir(line)
			if dir == "." {
				dir = "./"
			} else {
				dir = "./" + dir
			}
			if _, ok := seen[dir]; !ok {
				seen[dir] = struct{}{}
				dirs = append(dirs, dir)
			}
		}
	}

	return dirs, nil
}

func printMineSuccess() error {
	if err := mempalace.Status(); err != nil {
		return err
	}
	fmt.Println()
	fmt.Printf("  %s Knowledge graph updated!\n", color.GreenString("✓"))
	fmt.Printf("  %s Claude now knows about your new code\n", color.GreenString("✓"))
	fmt.Println()
	return nil
}

func pluralY(n int) string {
	if n == 1 {
		return "y"
	}
	return "ies"
}
