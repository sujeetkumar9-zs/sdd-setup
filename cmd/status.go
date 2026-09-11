package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

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

	fmt.Printf("  %s Mempalace Knowledge Graph:\n\n", color.BlueString("→"))
	mempalace.Status()

	fmt.Println()
	showCoverage()

	fmt.Println()
	fmt.Printf("  %s MCP Servers:\n\n", color.BlueString("→"))
	showMCPStatus()

	fmt.Println()
	return nil
}

func showMCPStatus() {
	out, err := exec.Command("claude", "mcp", "list").Output()
	if err != nil {
		fmt.Printf("  %s Could not run 'claude mcp list': %s\n",
			color.RedString("✗"), err.Error())
		return
	}
	fmt.Print(string(out))
}

// showCoverage prints a coverage line: indexed drawers vs source files.
func showCoverage() {
	drawers := parseDrawerCount()
	files := countSourceFiles(".")

	fmt.Printf("  %s Coverage:\n", color.BlueString("→"))

	if drawers < 0 || files == 0 {
		fmt.Printf("      %s\n\n", color.HiBlackString("(unable to calculate — run 'sdd mine' to index)"))
		return
	}

	pct := min(drawers*100/files, 100)

	bar := coverageBar(pct)
	pctStr := fmt.Sprintf("%d%%", pct)

	var pctColored string
	switch {
	case pct >= 80:
		pctColored = color.GreenString(pctStr)
	case pct >= 40:
		pctColored = color.YellowString(pctStr)
	default:
		pctColored = color.RedString(pctStr)
	}

	fmt.Printf("      %s %s  (%d drawers / %d source files)\n\n",
		bar, pctColored, drawers, files)
}

// parseDrawerCount runs mempalace status and extracts the drawer count.
func parseDrawerCount() int {
	out, err := exec.Command("mempalace", "--palace", palacePathForStatus(), "status").Output()
	if err != nil {
		return -1
	}
	// Look for a number followed by "drawer" anywhere in the output.
	re := regexp.MustCompile(`(\d+)\s+drawer`)
	m := re.FindSubmatch(out)
	if m == nil {
		// Fallback: look for "Drawers: N" pattern.
		re2 := regexp.MustCompile(`(?i)drawers?[:\s]+(\d+)`)
		m = re2.FindSubmatch(out)
	}
	if m == nil {
		return -1
	}
	n, _ := strconv.Atoi(string(m[1]))
	return n
}

// palacePathForStatus returns the absolute palace path (mirrors mempalace.palacePath).
func palacePathForStatus() string {
	abs, err := filepath.Abs(".mempalace/palace")
	if err != nil {
		return ".mempalace/palace"
	}
	return abs
}

// countSourceFiles walks the project and counts non-ignored source files.
func countSourceFiles(root string) int {
	skipDirs := map[string]bool{
		".git": true, ".venv": true, ".mempalace": true,
		"node_modules": true, "vendor": true,
		"dist": true, "build": true, "target": true,
		".claude": true,
	}
	srcExts := map[string]bool{
		".go": true, ".ts": true, ".js": true, ".tsx": true, ".jsx": true,
		".py": true, ".rs": true, ".java": true, ".kt": true,
		".rb": true, ".php": true, ".cs": true, ".cpp": true, ".c": true, ".h": true,
	}

	count := 0
	filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			if skipDirs[d.Name()] {
				return filepath.SkipDir
			}
			return nil
		}
		if srcExts[strings.ToLower(filepath.Ext(d.Name()))] {
			count++
		}
		return nil
	})
	return count
}

// coverageBar returns a simple ASCII progress bar for the given percentage.
func coverageBar(pct int) string {
	const width = 20
	filled := pct * width / 100
	bar := "[" + strings.Repeat("█", filled) + strings.Repeat("░", width-filled) + "]"
	return color.HiBlackString(bar)
}
