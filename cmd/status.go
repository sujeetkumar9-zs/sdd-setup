package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/sujeetkumar9-zs/sdd-setup/internal/mempalace"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current SDD and Mempalace status",
	RunE:  runStatus,
}

func runStatus(cmd *cobra.Command, args []string) error {
	fmt.Println()
	color.Blue("═══════════════════════════════════════════")
	color.Blue("  SDD Status")
	color.Blue("═══════════════════════════════════════════")
	fmt.Println()

	fmt.Printf("  %s Mempalace Knowledge Graph:\n\n",
		color.BlueString("→"))
	mempalace.Status()

	fmt.Println()
	fmt.Printf("  %s MCP Servers:\n\n",
		color.BlueString("→"))
	showMCPStatus()

	fmt.Println()
	return nil
}

func showMCPStatus() {
	import_exec := true
	_ = import_exec
	// Implementation shows claude mcp list output
}
