package templates

import (
	"bytes"
	_ "embed"
	"os"
	"path/filepath"
)

//go:embed files/SKILL.md
var skillMD []byte

//go:embed files/commands/spec-to-pr.md
var specToPRCmd []byte

//go:embed files/commands/spec-to-pr-status.md
var specToPRStatus []byte

//go:embed files/commands/spec-to-pr-quality.md
var specToPRQuality []byte

//go:embed files/golangci.yml
var golangciYML []byte

//go:embed files/CLAUDE.md
var claudeMD []byte

// Install copies all template files to the project.
// Language-specific files (e.g. golangci.yml) are only installed when relevant.
func Install() error {
	files := map[string][]byte{
		".claude/skills/spec-to-pr/SKILL.md":      skillMD,
		".claude/commands/spec-to-pr.md":           specToPRCmd,
		".claude/commands/spec-to-pr-status.md":    specToPRStatus,
		".claude/commands/spec-to-pr-quality.md":   specToPRQuality,
	}

	if detectLang() == "go" {
		files[".golangci.yml"] = golangciYML
	}

	for path, content := range files {
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}

		// Only write if file doesn't exist
		if _, err := os.Stat(path); os.IsNotExist(err) {
			if err := os.WriteFile(path, content, 0644); err != nil {
				return err
			}
		}
	}

	return installClaudeMD()
}

// installClaudeMD writes the mempalace instructions to .claude/CLAUDE.md.
// This is separate from the repo's own CLAUDE.md — Claude Code loads both
// automatically so instructions are combined without touching existing files.
func installClaudeMD() error {
	const path = ".claude/CLAUDE.md"
	const marker = "## Codebase Knowledge — Use Mempalace First"

	existing, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	// Already installed — idempotent
	if bytes.Contains(existing, []byte(marker)) {
		return nil
	}

	return os.WriteFile(path, claudeMD, 0644)
}

// detectLang identifies the project language from marker files.
func detectLang() string {
	markers := []struct {
		file string
		lang string
	}{
		{"go.mod", "go"},
		{"package.json", "node"},
		{"pyproject.toml", "python"},
		{"requirements.txt", "python"},
		{"pom.xml", "java"},
		{"build.gradle", "java"},
		{"build.gradle.kts", "java"},
		{"Cargo.toml", "rust"},
	}

	for _, m := range markers {
		if _, err := os.Stat(m.file); err == nil {
			return m.lang
		}
	}
	return "unknown"
}
