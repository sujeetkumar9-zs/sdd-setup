package templates

import (
	_ "embed"
	"os"
	"path/filepath"
)

// Embed all template files into the binary
// This means NO external files needed!

//go:embed files/SKILL.md
var skillMD []byte

//go:embed files/commands/spec-to-pr.md
var specToPRCmd []byte

//go:embed files/commands/spec-to-pr-status.md
var specToPRStatus []byte

//go:embed files/commands/spec-to-pr-quality.md
var specToPRQuality []byte

//go:embed files/examples/GOOD_SPEC.md
var goodSpec []byte

//go:embed files/examples/BAD_SPEC.md
var badSpec []byte

//go:embed files/golangci.yml
var golangciYML []byte

// Install copies all template files to the project
func Install() error {
	files := map[string][]byte{
		".claude/skills/spec-to-pr-go/SKILL.md":  skillMD,
		".claude/commands/spec-to-pr.md":         specToPRCmd,
		".claude/commands/spec-to-pr-status.md":  specToPRStatus,
		".claude/commands/spec-to-pr-quality.md": specToPRQuality,
		"examples/GOOD_SPEC.md":                  goodSpec,
		"examples/BAD_SPEC.md":                   badSpec,
		".golangci.yml":                          golangciYML,
	}

	for path, content := range files {
		// Create directory
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

	return nil
}
