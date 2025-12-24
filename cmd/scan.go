package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

type ProjectMatcher struct {
	Name       string
	TargetDir  string
	AnchorFile string
}

var matchers = []ProjectMatcher{
	{"Node.js", "node_modules", "package.json"},
	{"Rust", "target", "Cargo.toml"},
	{"Python", "venv", "requirements.txt"},
	{"Python", ".venv", "pyproject.toml"},
	{"Maven", "target", "pom.xml"},
	{"Gradle", "build", "build.gradle"},
}

// dirSize calculates the total size of a directory in bytes
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

// formatSize converts bytes to a human-readable string
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

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan for stale project dependencies",
	Run: func(cmd *cobra.Command, args []string) {
		path, _ := cmd.Flags().GetString("path")
		days, _ := cmd.Flags().GetInt("days")
		threshold := time.Duration(days) * 24 * time.Hour

		fmt.Printf("🔍 Scanning %s for projects older than %d days...\n\n", path, days)

		var totalSaved int64

		err := filepath.WalkDir(path, func(p string, d os.DirEntry, err error) error {
			if err != nil {
				return filepath.SkipDir
			}

			if !isSafe(p) { return filepath.SkipDir }

			if !d.IsDir() {
				return nil
			}

			// Skip hidden directories like .git
			if d.Name() == ".git" || d.Name() == ".cache" {
				return filepath.SkipDir
			}

			// Check matchers
			for _, m := range matchers {
				if d.Name() == m.TargetDir {
					parent := filepath.Dir(p)
					anchorPath := filepath.Join(parent, m.AnchorFile)

					// Check if Anchor File exists
					if info, err := os.Stat(anchorPath); err == nil {
						// Age Check
						if time.Since(info.ModTime()) > threshold {
							size, _ := dirSize(p)
							totalSaved += size
							fmt.Printf("[STALE] %-10s | %-10s | %s\n",
								m.Name, formatSize(size), p)
						}
						return filepath.SkipDir
					}
				}
			}
			return nil
		})

		if err != nil {
			fmt.Printf("\nError during scan: %v\n", err)
		}

		fmt.Printf("\nDone! Potential space to reclaim: %s\n", formatSize(totalSaved))
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.Flags().StringP("path", "p", ".", "Path to scan")
	scanCmd.Flags().IntP("days", "d", 30, "Threshold of days since last use")
}