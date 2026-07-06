package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "sdd",
	Version: "1.0.0",
	Short:   "SDD - Spec Driven Development Toolkit",
	Long: color.BlueString(`
╔═══════════════════════════════════════════════╗
║     SDD - Spec Driven Development Toolkit     ║
║     Kroger Technology | Go Projects           ║
╚═══════════════════════════════════════════════╝

Transforms Jira/Confluence specs into merged PRs
using Claude Code + Mempalace knowledge graph.
    `),
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(setupCmd)
	rootCmd.AddCommand(mineCmd)
	rootCmd.AddCommand(verifyCmd)
	rootCmd.AddCommand(qualityCmd)
	rootCmd.AddCommand(statusCmd)
}
