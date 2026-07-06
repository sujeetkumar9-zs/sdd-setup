package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var qualityCmd = &cobra.Command{
	Use:   "quality",
	Short: "Run all quality gates before creating a PR",
	Long: `Runs all quality checks required before creating a PR:

  1. go test ./...         (all tests must pass)
  2. go vet ./...          (static analysis)
  3. golangci-lint run     (lint checks)
  4. gofmt -l .            (format check)
  5. go build ./...        (must compile)

Blocks on any failure. Fix all issues before PR.`,
	RunE: runQuality,
}

type gate struct {
	name    string
	command []string
	check   func(output string) error
}

func runQuality(cmd *cobra.Command, args []string) error {
	fmt.Println()
	color.Blue("═══════════════════════════════════════════")
	color.Blue("  SDD Quality Gates")
	color.Blue("═══════════════════════════════════════════")
	fmt.Println()

	gates := []gate{
		{
			name:    "Tests",
			command: []string{"go", "test", "./..."},
		},
		{
			name:    "Race Condition Check",
			command: []string{"go", "test", "-race", "./..."},
		},
		{
			name:    "Vet",
			command: []string{"go", "vet", "./..."},
		},
		{
			name:    "Lint",
			command: []string{"golangci-lint", "run"},
		},
		{
			name:    "Format",
			command: []string{"gofmt", "-l", "."},
			check: func(output string) error {
				if output != "" {
					return fmt.Errorf(
						"files need formatting:\n%s\n  Run: gofmt -w .",
						output,
					)
				}
				return nil
			},
		},
		{
			name:    "Build",
			command: []string{"go", "build", "./..."},
		},
	}

	allPassed := true
	results := []string{}

	for _, g := range gates {
		fmt.Printf("  %s Running %s...\n",
			color.BlueString("→"), g.name)

		c := exec.Command(g.command[0], g.command[1:]...)
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		out, err := c.Output()

		if g.check != nil {
			err = g.check(string(out))
		}

		if err != nil {
			fmt.Printf("  %s %s FAILED\n\n",
				color.RedString("✗"), g.name)
			results = append(results,
				fmt.Sprintf("  %s %s", color.RedString("✗"), g.name))
			allPassed = false
		} else {
			fmt.Printf("  %s %s passed\n\n",
				color.GreenString("✓"), g.name)
			results = append(results,
				fmt.Sprintf("  %s %s", color.GreenString("✓"), g.name))
		}
	}

	fmt.Println()
	color.Blue("═══════════════════════════════════════════")
	fmt.Println("  Summary:")
	fmt.Println()
	for _, r := range results {
		fmt.Println(r)
	}
	fmt.Println()

	if allPassed {
		color.Green("  ✅ All quality gates passed! Ready for PR.")
	} else {
		color.Red("  ❌ Quality gates failed. Fix before creating PR.")
		return fmt.Errorf("quality gates failed")
	}

	fmt.Println()
	return nil
}
