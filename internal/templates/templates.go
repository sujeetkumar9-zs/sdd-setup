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

//go:embed files/commands/spec-to-pr-resume.md
var specToPRResume []byte

//go:embed files/commands/spec-to-pr-quality.md
var specToPRQuality []byte

//go:embed files/commands/review-pr.md
var reviewPRCmd []byte

//go:embed files/skills/review-pr/SKILL.md
var reviewPRSkill []byte

//go:embed files/commands/address-pr-comments.md
var addressPRCommentsCmd []byte

//go:embed files/skills/address-pr-comments/SKILL.md
var addressPRCommentsSkill []byte

//go:embed files/commands/fix-bug.md
var fixBugCmd []byte

//go:embed files/skills/fix-bug/SKILL.md
var fixBugSkill []byte

//go:embed files/golangci.yml
var golangciYML []byte

//go:embed files/CLAUDE.md
var claudeMD []byte

// Install copies all template files to the project, skipping files that already exist.
// Language-specific files (e.g. golangci.yml) are only installed when relevant.
func Install() error {
	return install(false)
}

// InstallForce overwrites all template files unconditionally.
// Used by `sdd update` to refresh skills and commands to the latest version.
func InstallForce() (int, error) {
	return installCount(true)
}

func install(force bool) error {
	_, err := installCount(force)
	return err
}

func installCount(force bool) (int, error) {
	files := map[string][]byte{
		".claude/skills/spec-to-pr/SKILL.md":           skillMD,
		".claude/commands/spec-to-pr.md":                specToPRCmd,
		".claude/commands/spec-to-pr-status.md":         specToPRStatus,
		".claude/commands/spec-to-pr-resume.md":         specToPRResume,
		".claude/commands/spec-to-pr-quality.md":        specToPRQuality,
		".claude/commands/review-pr.md":                 reviewPRCmd,
		".claude/skills/review-pr/SKILL.md":             reviewPRSkill,
		".claude/commands/address-pr-comments.md":       addressPRCommentsCmd,
		".claude/skills/address-pr-comments/SKILL.md":   addressPRCommentsSkill,
		".claude/commands/fix-bug.md":                   fixBugCmd,
		".claude/skills/fix-bug/SKILL.md":               fixBugSkill,
	}

	if detectLang() == "go" {
		files[".golangci.yml"] = golangciYML
	}

	written := 0
	for path, content := range files {
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return written, err
		}

		if !force {
			if _, err := os.Stat(path); !os.IsNotExist(err) {
				continue
			}
		}

		if err := os.WriteFile(path, content, 0644); err != nil {
			return written, err
		}
		written++
	}

	if err := installClaudeMD(force); err != nil {
		return written, err
	}
	return written, nil
}

// installClaudeMD writes the mempalace instructions to .claude/CLAUDE.md.
// This is separate from the repo's own CLAUDE.md — Claude Code loads both
// automatically so instructions are combined without touching existing files.
func installClaudeMD(force bool) error {
	const path = ".claude/CLAUDE.md"
	const marker = "## Codebase Knowledge — Use Mempalace First"

	if !force {
		existing, err := os.ReadFile(path)
		if err != nil && !os.IsNotExist(err) {
			return err
		}
		// Already installed — idempotent
		if bytes.Contains(existing, []byte(marker)) {
			return nil
		}
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
