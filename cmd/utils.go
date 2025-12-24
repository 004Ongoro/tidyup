package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// ProjectMatcher defines the structure for identifying a project type
type ProjectMatcher struct {
	Name       string `mapstructure:"name"`
	TargetDir  string `mapstructure:"target_dir"`
	AnchorFile string `mapstructure:"anchor_file"`
}

// getBlocklist retrieves the list from config or returns defaults
func getBlocklist() []string {
	defaults := []string{
		"AppData", "Library", ".vscode", ".antigravity", ".rustup", ".cargo",
		"Program Files", "Windows", "System32", "node_modules",
	}

	if viper.IsSet("blocklist") {
		return viper.GetStringSlice("blocklist")
	}
	return defaults
}

// getMatchers retrieves matchers from config or returns defaults
func getMatchers() []ProjectMatcher {
	defaults := []ProjectMatcher{
		{"Node.js", "node_modules", "package.json"},
		{"Rust", "target", "Cargo.toml"},
		{"Python", "venv", "requirements.txt"},
		{"Python", ".venv", "pyproject.toml"},
		{"Maven", "target", "pom.xml"},
		{"Gradle", "build", "build.gradle"},
	}

	if viper.IsSet("matchers") {
		var custom []ProjectMatcher
		if err := viper.UnmarshalKey("matchers", &custom); err == nil {
			return custom
		}
	}
	return defaults
}

func isSafe(path string) bool {
	parts := strings.Split(path, string(os.PathSeparator))
	blocklist := getBlocklist()

	for _, part := range parts {
		for _, blocked := range blocklist {
			// Blocking if a PARENT folder is in the blocklist.
			if strings.EqualFold(part, blocked) && !isTargetDir(part) {
				return false
			}
		}
	}
	return true
}

func isTargetDir(name string) bool {
	matchers := getMatchers()
	for _, m := range matchers {
		if name == m.TargetDir {
			return true
		}
	}
	return false
}

func dirSize(path string) (int64, error) {
	var size int64
	err := filepath.WalkDir(path, func(_ string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			info, err := d.Info()
			if err != nil {
				return err
			}
			size += info.Size()
		}
		return nil
	})
	return size, err
}

func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}