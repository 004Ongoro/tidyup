package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

// Blocklist of directories that TidyUp should NEVER touch
var blocklist = []string{
	"AppData", "Library", ".vscode", ".antigravity", ".rustup", ".cargo", 
    "Program Files", "Windows", "System32",
}

func isSafe(path string) bool {
	// Split the path into individual directory names
	parts := strings.Split(path, string(os.PathSeparator))
	
	for _, part := range parts {
		for _, blocked := range blocklist {
			// Exact match for a blocked directory name
			if strings.EqualFold(part, blocked) {
				return false
			}
		}
	}
	return true
}

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Find and delete stale dependency folders",
	Run: func(cmd *cobra.Command, args []string) {
		path, _ := cmd.Flags().GetString("path")
		days, _ := cmd.Flags().GetInt("days")
		force, _ := cmd.Flags().GetBool("force")
		threshold := time.Duration(days) * 24 * time.Hour

		fmt.Printf("🧹 TidyUp is looking for targets in: %s\n", path)

		var targets []string
		err := filepath.WalkDir(path, func(p string, d os.DirEntry, err error) error {
			if err != nil || !d.IsDir() || !isSafe(p) {
				return nil
			}

			for _, m := range matchers {
				if d.Name() == m.TargetDir {
					anchorPath := filepath.Join(filepath.Dir(p), m.AnchorFile)
					if info, err := os.Stat(anchorPath); err == nil {
						if time.Since(info.ModTime()) > threshold {
							targets = append(targets, p)
						}
						return filepath.SkipDir
					}
				}
			}
			return nil
		})

		if err != nil || len(targets) == 0 {
			fmt.Println("✨ No stale folders found. Your machine is already tidy!")
			return
		}

		var toDelete []string

		if force {
			toDelete = targets
			fmt.Printf("⚠️ Force mode enabled. Deleting %d folders...\n", len(targets))
		} else {
			// Interactive Mode
			prompt := &survey.MultiSelect{
				Message: "Select the folders you want to delete:",
				Options: targets,
			}
			survey.AskOne(prompt, &toDelete)
		}

		for _, folder := range toDelete {
			fmt.Printf("Deleting: %s...", folder)
			err := os.RemoveAll(folder)
			if err != nil {
				fmt.Printf(" ❌ Error: %v\n", err)
			} else {
				fmt.Println(" ✅ Done")
			}
		}

		fmt.Printf("\n⭐ Cleanup complete! Deleted %d folders.\n", len(toDelete))
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
	cleanCmd.Flags().StringP("path", "p", ".", "Path to scan")
	cleanCmd.Flags().IntP("days", "d", 30, "Threshold of days")
	cleanCmd.Flags().BoolP("force", "f", false, "Delete everything without asking")
}