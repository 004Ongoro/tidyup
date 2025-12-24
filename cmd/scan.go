package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type ScanResult struct {
	Type string
	Path string
	Size int64
	Time time.Time
}

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan for stale project dependencies",
	Long: `Scans the provided directory for known project types. 
It checks the modification date of anchor files to determine if a project is stale.`,
	Run: func(cmd *cobra.Command, args []string) {
		path, _ := cmd.Flags().GetString("path")
		days, _ := cmd.Flags().GetInt("days")
		threshold := time.Duration(days) * 24 * time.Hour

		matchers := getMatchers()

		color.Cyan("TidyUp Scanning: %s (Older than %d days)\n", path, days)

		results := make(chan ScanResult)
		var wg sync.WaitGroup
		var totalSaved int64
		count := 0

		stopSpinner := make(chan bool)
		go func() {
			for {
				select {
				case <-stopSpinner:
					return
				default:
					fmt.Print(".")
					time.Sleep(500 * time.Millisecond)
				}
			}
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			filepath.WalkDir(path, func(p string, d os.DirEntry, err error) error {
				if err != nil || !isSafe(p) {
					return filepath.SkipDir
				}
				if !d.IsDir() || d.Name() == ".git" {
					return nil
				}

				for _, m := range matchers {
					if d.Name() == m.TargetDir {
						anchor := filepath.Join(filepath.Dir(p), m.AnchorFile)
						if info, err := os.Stat(anchor); err == nil {
							if time.Since(info.ModTime()) > threshold {
								size, _ := dirSize(p)
								results <- ScanResult{m.Name, p, size, info.ModTime()}
							}
							return filepath.SkipDir
						}
					}
				}
				return nil
			})
			close(results)
		}()

		fmt.Println()
		for res := range results {
			count++
			totalSaved += res.Size
			color.Red("[STALE] ")
			fmt.Printf("%-10s ", res.Type)
			color.Green("%-10s ", formatSize(res.Size))
			color.HiBlack("| %s (%s)\n", res.Path, res.Time.Format("2006-01-02"))
		}

		stopSpinner <- true
		wg.Wait()
		fmt.Println("\n" + strings.Repeat("-", 60))
		color.HiYellow("Summary: Found %d folders | Total Space: %s", count, formatSize(totalSaved))
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.Flags().StringP("path", "p", ".", "Path to scan")
	scanCmd.Flags().IntP("days", "d", 30, "Age threshold")
}