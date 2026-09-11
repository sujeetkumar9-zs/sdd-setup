package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var qualityCmd = &cobra.Command{
	Use:   "quality",
	Short: "Run all quality gates before creating a PR",
	Long: `Detects the project language and runs the appropriate quality checks:

  Go:     go test, go vet, golangci-lint, gofmt, go build
  Node:   npm test, npm run lint, npm run build
  Python: pytest, ruff check, ruff format --check
  Rust:   cargo test, cargo clippy, cargo fmt --check, cargo build
  Java:   mvn/gradlew test + build

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

	gates := languageGates()

	allPassed := true
	results := []string{}

	for _, g := range gates {
		fmt.Printf("  %s Running %s...\n",
			color.BlueString("→"), g.name)

		c := exec.Command(g.command[0], g.command[1:]...)

		var buf bytes.Buffer
		if g.check != nil {
			// Gates with a check function need captured output — buffer only.
			c.Stdout = &buf
			c.Stderr = &buf
		} else {
			// Stream live output AND capture for RCA on failure.
			c.Stdout = io.MultiWriter(os.Stdout, &buf)
			c.Stderr = io.MultiWriter(os.Stderr, &buf)
		}

		runErr := c.Run()

		var err error
		if g.check != nil {
			err = g.check(buf.String())
		} else {
			err = runErr
		}

		if err != nil {
			if g.check != nil && buf.Len() > 0 {
				// For check-based gates the buffered output wasn't streamed — print it now.
				fmt.Print(buf.String())
			}
			fmt.Printf("  %s %s FAILED: %s\n\n",
				color.RedString("✗"), g.name, color.RedString(err.Error()))
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

// languageGates detects the project language and returns the appropriate gates.
func languageGates() []gate {
	for _, m := range []struct {
		file  string
		gates []gate
	}{
		{"go.mod", goGates()},
		{"package.json", nodeGates()},
		{"pyproject.toml", pythonGates()},
		{"requirements.txt", pythonGates()},
		{"Cargo.toml", rustGates()},
		{"pom.xml", mavenGates()},
		{"build.gradle", gradleGates()},
		{"build.gradle.kts", gradleGates()},
	} {
		if _, err := os.Stat(m.file); err == nil {
			return m.gates
		}
	}
	return goGates()
}

func goGates() []gate {
	return []gate{
		{name: "Tests", command: []string{"go", "test", "./..."}},
		{name: "Race Condition Check", command: []string{"go", "test", "-race", "./..."}},
		{name: "Vet", command: []string{"go", "vet", "./..."}},
		{name: "Lint", command: []string{"golangci-lint", "run"}},
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
		{name: "Build", command: []string{"go", "build", "./..."}},
	}
}

func nodeGates() []gate {
	return []gate{
		{name: "Tests", command: []string{"npm", "test"}},
		{name: "Lint", command: []string{"npm", "run", "lint"}},
		{name: "Build", command: []string{"npm", "run", "build"}},
	}
}

func pythonGates() []gate {
	return []gate{
		{name: "Tests", command: []string{"pytest"}},
		{name: "Lint", command: []string{"ruff", "check", "."}},
		{
			name:    "Format",
			command: []string{"ruff", "format", "--check", "."},
		},
	}
}

func rustGates() []gate {
	return []gate{
		{name: "Tests", command: []string{"cargo", "test"}},
		{name: "Lint", command: []string{"cargo", "clippy", "--", "-D", "warnings"}},
		{
			name:    "Format",
			command: []string{"cargo", "fmt", "--check"},
		},
		{name: "Build", command: []string{"cargo", "build"}},
	}
}

func mavenGates() []gate {
	return []gate{
		{name: "Tests", command: []string{"mvn", "test"}},
		{name: "Build", command: []string{"mvn", "package", "-DskipTests"}},
	}
}

func gradleGates() []gate {
	return []gate{
		{name: "Tests", command: []string{"./gradlew", "test"}},
		{name: "Build", command: []string{"./gradlew", "build", "-x", "test"}},
	}
}
