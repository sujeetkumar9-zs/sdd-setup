package system

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)


// Language constants returned by DetectLanguage.
const (
	LangGo      = "go"
	LangNode    = "node"
	LangPython  = "python"
	LangJava    = "java"
	LangRust    = "rust"
	LangUnknown = "unknown"
)

// DetectLanguage inspects the current directory and returns the primary language.
func DetectLanguage() string {
	markers := []struct {
		file string
		lang string
	}{
		{"go.mod", LangGo},
		{"package.json", LangNode},
		{"pyproject.toml", LangPython},
		{"requirements.txt", LangPython},
		{"pom.xml", LangJava},
		{"build.gradle", LangJava},
		{"build.gradle.kts", LangJava},
		{"Cargo.toml", LangRust},
	}

	for _, m := range markers {
		if _, err := os.Stat(m.file); err == nil {
			return m.lang
		}
	}
	return LangUnknown
}

// CheckPrerequisites verifies all required tools are installed for the given language.
func CheckPrerequisites(lang string) error {
	required := map[string]string{
		"python3":   "Run: brew install python",
		"claude":    "Run: npm install -g @anthropic-ai/claude-code",
		"mempalace": "Run: pipx install mempalace",
	}

	switch lang {
	case LangGo:
		required["go"] = "Install from: https://go.dev/dl/"
	case LangNode:
		required["node"] = "Install from: https://nodejs.org"
		required["npm"] = "Install from: https://nodejs.org"
	case LangRust:
		required["cargo"] = "Install from: https://rustup.rs"
	case LangJava:
		// Maven or Gradle — check whichever is present
		if _, err := exec.LookPath("mvn"); err != nil {
			if _, err2 := exec.LookPath("gradle"); err2 != nil {
				required["mvn"] = "Install Maven: https://maven.apache.org or Gradle: https://gradle.org"
			}
		}
	}

	var missing []string
	for tool, hint := range required {
		if _, err := exec.LookPath(tool); err != nil {
			missing = append(missing, fmt.Sprintf("%s (%s)", tool, hint))
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing tools:\n  • %s", strings.Join(missing, "\n  • "))
	}
	return nil
}

// CreateVenv creates a Python virtual environment for mempalace.
func CreateVenv() error {
	if _, err := os.Stat(".venv"); os.IsNotExist(err) {
		out, err := exec.Command("python3", "-m", "venv", ".venv").CombinedOutput()
		if err != nil {
			return fmt.Errorf("python3 -m venv .venv failed\n  RCA: %w\n  Output: %s", err, string(out))
		}
	}
	return nil
}

// InstallLanguageTools installs quality tools appropriate for the detected language.
func InstallLanguageTools(lang string) error {
	switch lang {
	case LangGo:
		return installGoTools()
	case LangNode:
		return installNodeTools()
	case LangPython:
		return installPythonTools()
	default:
		// Rust: rustup manages clippy/rustfmt — nothing to install
		// Java: Maven/Gradle manage tools — nothing to install
		// Unknown: skip
		return nil
	}
}

func installGoTools() error {
	tools := []string{
		"github.com/golangci/golangci-lint/cmd/golangci-lint@latest",
		"github.com/vektra/mockery/v2@latest",
	}

	for _, tool := range tools {
		out, err := exec.Command("go", "install", tool).CombinedOutput()
		if err != nil {
			return fmt.Errorf("go install %s failed\n  RCA: %w\n  Output: %s", tool, err, string(out))
		}
	}
	return nil
}

func installNodeTools() error {
	if _, err := exec.LookPath("eslint"); err != nil {
		out, err := exec.Command("npm", "install", "-g", "eslint").CombinedOutput()
		if err != nil {
			return fmt.Errorf("npm install -g eslint failed\n  RCA: %w\n  Output: %s", err, string(out))
		}
	}
	return nil
}

func installPythonTools() error {
	tools := []string{"ruff", "pytest"}
	for _, tool := range tools {
		if _, err := exec.LookPath(tool); err != nil {
			out, err := exec.Command("pip", "install", tool).CombinedOutput()
			if err != nil {
				return fmt.Errorf("pip install %s failed\n  RCA: %w\n  Output: %s", tool, err, string(out))
			}
		}
	}
	return nil
}

// UpdateGitignore adds required entries to .gitignore.
func UpdateGitignore() error {
	content, _ := os.ReadFile(".gitignore")
	existing := string(content)

	additions := "\n# SDD - Spec Driven Development\n"
	changed := false

	entries := []string{".mempalace/", ".venv/", "__pycache__/"}
	for _, entry := range entries {
		if !strings.Contains(existing, entry) {
			additions += entry + "\n"
			changed = true
		}
	}

	if changed {
		f, err := os.OpenFile(".gitignore",
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = f.WriteString(additions)
		return err
	}
	return nil
}
