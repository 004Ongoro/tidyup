package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Interactively delete stale dependencies",
	Run: func(cmd *cobra.Command, args []string) {
		path, _ := cmd.Flags().GetString("path")
		days, _ := cmd.Flags().GetInt("days")
		force, _ := cmd.Flags().GetBool("force")
		deep, _ := cmd.Flags().GetBool("deep")
		threshold := time.Duration(days) * 24 * time.Hour

		matchers := getMatchers()

		color.Magenta(" TidyUp Cleanup initialized...")
		if deep {
			color.Red("  Deep Mode: Targeting all matched folders regardless of anchor files.")
		}

		var targets []string
		var scannedCount int

		fmt.Print("🔍 Looking for targets...")
		filepath.WalkDir(path, func(p string, d os.DirEntry, err error) error {
			scannedCount++
			if err != nil || !isSafe(p) {
				return filepath.SkipDir
			}
			if !d.IsDir() {
				return nil
			}

			for _, m := range matchers {
				if d.Name() == m.TargetDir {
					if deep {
						info, err := d.Info()
						if err == nil && time.Since(info.ModTime()) > threshold {
							targets = append(targets, p)
						}
						return filepath.SkipDir
					} else {
						anchor := filepath.Join(filepath.Dir(p), m.AnchorFile)
						if info, err := os.Stat(anchor); err == nil {
							if time.Since(info.ModTime()) > threshold {
								targets = append(targets, p)
							}
							return filepath.SkipDir
						}
					}
				}
			}
			return nil
		})
		fmt.Printf(" Done! (Scanned %d directories)\n", scannedCount)

		if len(targets) == 0 {
			color.Cyan(" Nothing to clean! Your drive is tidy.")
			return
		}

		var toDelete []string
		if force {
			toDelete = targets
		} else {
			prompt := &survey.MultiSelect{
				Message: "Select folders to PERMANENTLY delete:",
				Options: targets,
			}
			survey.AskOne(prompt, &toDelete)
		}

		for _, folder := range toDelete {
			fmt.Printf("Removing %s...", folder)
			if err := os.RemoveAll(folder); err != nil {
				color.Red("  Error: %v", err)
			} else {
				color.Green("  Done")
			}
		}

		color.HiMagenta("\n Finished! Cleaned %d folders.", len(toDelete))
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
	cleanCmd.Flags().StringP("path", "p", ".", "Path to scan")
	cleanCmd.Flags().IntP("days", "d", 30, "Age threshold")
	cleanCmd.Flags().BoolP("force", "f", false, "Skip confirmation")
	cleanCmd.Flags().Bool("deep", false, "Perform deep scan (ignore anchor files)")
}
