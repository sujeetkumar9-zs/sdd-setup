package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// Version is set at build time via -ldflags "-X github.com/sujeetkumar9-zs/sdd-setup/cmd.Version=..."
var Version = "dev"

var rootCmd = &cobra.Command{
	Use:     "sdd",
	Version: Version,
	Short:   "SDD - Spec Driven Development Toolkit",
	Long: color.BlueString(`
╔═══════════════════════════════════════════════╗
║     SDD - Spec Driven Development Toolkit     ║
║     Spec-Driven Development for Go            ║
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
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(teardownCmd)
}
