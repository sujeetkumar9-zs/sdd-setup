package cmd

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/sujeetkumar9-zs/sdd-setup/internal/mempalace"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search the project knowledge graph",
	Long: `Searches the project-scoped Mempalace knowledge graph.

Always scoped to the current project (.mempalace/palace).
Use this instead of 'mempalace search' to avoid seeing
results from other projects.

Example:
  sdd search "explain the household card endpoint flow"`,
	Args: cobra.MinimumNArgs(1),
	RunE: runSearch,
}

func runSearch(cmd *cobra.Command, args []string) error {
	query := strings.Join(args, " ")

	fmt.Println()
	fmt.Printf("  %s Searching knowledge graph...\n", color.BlueString("→"))
	fmt.Println()

	if err := mempalace.Search(query); err != nil {
		return fmt.Errorf("search failed: %w", err)
	}

	fmt.Println()
	return nil
}
